//! eg — EventGraph CLI
//!
//! Usage:
//!   eg bootstrap              Initialize a new event graph
//!   eg get <event-id>         Get an event by ID
//!   eg recent [limit]         Show recent events (default: 10)
//!   eg count                  Show event count
//!   eg verify                 Verify hash chain integrity
//!   eg head                   Show the chain head
//!   eg help                   Show this help

use std::env;
use std::process;

use eventgraph::actor::InMemoryActorStore;
use eventgraph::graph::Graph;
use eventgraph::store::{InMemoryStore, Store};
use eventgraph::types::{ActorId, EventId};

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        usage();
        process::exit(1);
    }

    let cmd = &args[1];
    let store = InMemoryStore::new();
    let actor_store = InMemoryActorStore::new();
    let mut g = Graph::new(store, actor_store);
    g.start().unwrap_or_else(|e| fatal(&format!("graph start failed: {e}")));

    match cmd.as_str() {
        "bootstrap" => {
            let actor_id = ActorId::new("actor_system0000000000000000000001")
                .unwrap_or_else(|e| fatal(&format!("invalid actor ID: {e}")));
            match g.bootstrap(actor_id, None) {
                Ok(ev) => print_event(&ev),
                Err(e) => fatal(&format!("bootstrap failed: {e}")),
            }
        }
        "get" => {
            if args.len() < 3 {
                fatal("usage: eg get <event-id>");
            }
            let id = EventId::new(&args[2])
                .unwrap_or_else(|e| fatal(&format!("invalid event ID: {e}")));
            match g.store().get(&id) {
                Ok(ev) => print_event(ev),
                Err(e) => fatal(&format!("get failed: {e}")),
            }
        }
        "recent" => {
            let limit = if args.len() >= 3 {
                args[2].parse::<usize>().unwrap_or_else(|_| {
                    fatal(&format!("invalid limit: {}", args[2]));
                })
            } else {
                10
            };
            let events = g.store().recent(limit);
            if events.is_empty() {
                println!("(no events)");
            } else {
                for ev in &events {
                    print_event_summary(ev);
                }
            }
        }
        "count" => {
            println!("{} events", g.store().count());
        }
        "verify" => {
            let result = g.store().verify_chain();
            println!(
                "Chain verified: {} events, valid={}",
                result.length, result.valid
            );
        }
        "head" => match g.store().head() {
            Some(ev) => print_event(ev),
            None => println!("(empty chain)"),
        },
        "help" | "-h" | "--help" => usage(),
        _ => {
            eprintln!("unknown command: {cmd}");
            usage();
            process::exit(1);
        }
    }

    g.close();
}

fn print_event(ev: &eventgraph::event::Event) {
    println!(
        r#"{{
  "id": "{}",
  "type": "{}",
  "source": "{}",
  "timestamp_nanos": {},
  "hash": "{}",
  "prev_hash": "{}",
  "conversation_id": "{}",
  "version": {}
}}"#,
        ev.id.value(),
        ev.event_type.value(),
        ev.source.value(),
        ev.timestamp_nanos,
        ev.hash.value(),
        ev.prev_hash.value(),
        ev.conversation_id.value(),
        ev.version,
    );
}

fn print_event_summary(ev: &eventgraph::event::Event) {
    println!(
        "  {}  {}  {}",
        ev.id.value(), ev.event_type.value(), ev.timestamp_nanos
    );
}

fn usage() {
    println!(
        "eg — EventGraph CLI

Usage:
  eg bootstrap              Initialize a new event graph
  eg get <event-id>         Get an event by ID
  eg recent [limit]         Show recent events (default: 10)
  eg count                  Show event count
  eg verify                 Verify hash chain integrity
  eg head                   Show the chain head
  eg help                   Show this help"
    );
}

fn fatal(msg: &str) -> ! {
    eprintln!("{msg}");
    process::exit(1);
}
