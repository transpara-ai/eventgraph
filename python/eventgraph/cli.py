#!/usr/bin/env python3
"""
eg — EventGraph CLI

Usage:
    eg bootstrap              Initialize a new event graph
    eg get <event-id>         Get an event by ID
    eg recent [limit]         Show recent events (default: 10)
    eg count                  Show event count
    eg verify                 Verify hash chain integrity
    eg head                   Show the chain head
    eg help                   Show this help
"""

import json
import sys
from datetime import datetime, timezone

from .event import Event, create_bootstrap, create_event
from .store import InMemoryStore
from .actor import InMemoryActorStore
from .graph import Graph
from .types import ActorID, EventID


def print_event(ev: Event) -> None:
    out = {
        "id": ev.id.value,
        "type": ev.type.value,
        "source": ev.source.value,
        "timestamp": datetime.fromtimestamp(
            ev.timestamp_nanos / 1_000_000_000, tz=timezone.utc
        ).isoformat(),
        "hash": ev.hash.value,
        "prev_hash": ev.prev_hash.value,
        "conversation_id": ev.conversation_id.value,
        "version": ev.version,
    }
    print(json.dumps(out, indent=2))


def print_event_summary(ev: Event) -> None:
    ts = datetime.fromtimestamp(
        ev.timestamp_nanos / 1_000_000_000, tz=timezone.utc
    ).isoformat()
    print(f"  {ev.id.value}  {ev.type.value}  {ts}")


def usage() -> None:
    print("""eg — EventGraph CLI

Usage:
  eg bootstrap              Initialize a new event graph
  eg get <event-id>         Get an event by ID
  eg recent [limit]         Show recent events (default: 10)
  eg count                  Show event count
  eg verify                 Verify hash chain integrity
  eg head                   Show the chain head
  eg help                   Show this help""")


def fatal(msg: str) -> None:
    print(msg, file=sys.stderr)
    sys.exit(1)


def main() -> None:
    args = sys.argv[1:]
    if not args:
        usage()
        sys.exit(1)

    cmd = args[0]
    store = InMemoryStore()
    actor_store = InMemoryActorStore()
    g = Graph(store, actor_store)

    if cmd == "bootstrap":
        actor_id = ActorID("actor_system0000000000000000000001")
        ev = g.bootstrap(actor_id)
        print_event(ev)

    elif cmd == "get":
        if len(args) < 2:
            fatal("usage: eg get <event-id>")
        g.start()
        ev = store.get(EventID(args[1]))
        print_event(ev)

    elif cmd == "recent":
        g.start()
        limit = 10
        if len(args) >= 2:
            try:
                limit = int(args[1])
            except ValueError:
                fatal(f"invalid limit: {args[1]}")
        events = store.recent(limit)
        for ev in events:
            print_event_summary(ev)
        if not events:
            print("(no events)")

    elif cmd == "count":
        g.start()
        print(f"{store.count()} events")

    elif cmd == "verify":
        g.start()
        result = store.verify_chain()
        print(f"Chain verified: {result['length']} events, valid={result['valid']}")

    elif cmd == "head":
        g.start()
        head = store.head()
        if head is None:
            print("(empty chain)")
        else:
            print_event(head)

    elif cmd in ("help", "-h", "--help"):
        usage()

    else:
        print(f"unknown command: {cmd}", file=sys.stderr)
        usage()
        sys.exit(1)

    g.close()


if __name__ == "__main__":
    main()
