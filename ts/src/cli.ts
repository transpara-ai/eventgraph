#!/usr/bin/env node
/**
 * eg — EventGraph CLI
 *
 * Usage:
 *   eg bootstrap              Initialize a new event graph
 *   eg get <event-id>         Get an event by ID
 *   eg recent [limit]         Show recent events (default: 10)
 *   eg count                  Show event count
 *   eg verify                 Verify hash chain integrity
 *   eg head                   Show the chain head
 *   eg help                   Show this help
 */

import { ActorId, EventId } from "./types.js";
import { InMemoryStore } from "./store.js";
import { InMemoryActorStore } from "./actor.js";
import { Graph } from "./graph.js";
import { Event } from "./event.js";

function printEvent(ev: Event): void {
  const out = {
    id: ev.id.value,
    type: ev.type.value,
    source: ev.source.value,
    timestamp: new Date(ev.timestampNanos / 1_000_000).toISOString(),
    hash: ev.hash.value,
    prev_hash: ev.prevHash.value,
    conversation_id: ev.conversationId.value,
    version: ev.version,
  };
  console.log(JSON.stringify(out, null, 2));
}

function printEventSummary(ev: Event): void {
  const ts = new Date(ev.timestampNanos / 1_000_000).toISOString();
  console.log(`  ${ev.id.value}  ${ev.type.value}  ${ts}`);
}

function usage(): void {
  console.log(`eg — EventGraph CLI

Usage:
  eg bootstrap              Initialize a new event graph
  eg get <event-id>         Get an event by ID
  eg recent [limit]         Show recent events (default: 10)
  eg count                  Show event count
  eg verify                 Verify hash chain integrity
  eg head                   Show the chain head
  eg help                   Show this help`);
}

function fatal(msg: string): never {
  console.error(msg);
  process.exit(1);
}

const args = process.argv.slice(2);
if (args.length === 0) {
  usage();
  process.exit(1);
}

const cmd = args[0];
const store = new InMemoryStore();
const actorStore = new InMemoryActorStore();
const g = new Graph(store, actorStore);

switch (cmd) {
  case "bootstrap": {
    const actorId = new ActorId("actor_system0000000000000000000001");
    const ev = g.bootstrap(actorId);
    printEvent(ev);
    break;
  }
  case "get": {
    if (args.length < 2) fatal("usage: eg get <event-id>");
    g.start();
    const ev = store.get(new EventId(args[1]));
    printEvent(ev);
    break;
  }
  case "recent": {
    g.start();
    const limit = args.length >= 2 ? parseInt(args[1], 10) : 10;
    if (isNaN(limit)) fatal(`invalid limit: ${args[1]}`);
    const events = store.recent(limit);
    for (const ev of events) printEventSummary(ev);
    if (events.length === 0) console.log("(no events)");
    break;
  }
  case "count": {
    g.start();
    console.log(`${store.count()} events`);
    break;
  }
  case "verify": {
    g.start();
    const result = store.verifyChain();
    console.log(`Chain verified: ${result.length} events, valid=${result.valid}`);
    break;
  }
  case "head": {
    g.start();
    const head = store.head();
    if (head.isNone) {
      console.log("(empty chain)");
    } else {
      printEvent(head.unwrap());
    }
    break;
  }
  case "help":
  case "-h":
  case "--help":
    usage();
    break;
  default:
    console.error(`unknown command: ${cmd}`);
    usage();
    process.exit(1);
}

g.close();
