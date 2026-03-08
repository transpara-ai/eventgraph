// eg — EventGraph CLI
//
// Usage:
//   eg bootstrap              Initialize a new event graph
//   eg get <event-id>         Get an event by ID
//   eg recent [limit]         Show recent events (default: 10)
//   eg count                  Show event count
//   eg verify                 Verify hash chain integrity
//   eg head                   Show the chain head
//   eg help                   Show this help

using System.Text.Json;
using EventGraph;

if (args.Length == 0)
{
    Usage();
    return 1;
}

var cmd = args[0];
var store = new InMemoryStore();
var actorStore = new InMemoryActorStore();
var g = new Graph(store, actorStore);

try
{
    switch (cmd)
    {
        case "bootstrap":
        {
            var actorId = new ActorId("actor_system0000000000000000000001");
            var ev = g.Bootstrap(actorId);
            PrintEvent(ev);
            break;
        }
        case "get":
        {
            if (args.Length < 2) { Console.Error.WriteLine("usage: eg get <event-id>"); return 1; }
            g.Start();
            var ev = store.Get(new EventId(args[1]));
            PrintEvent(ev);
            break;
        }
        case "recent":
        {
            g.Start();
            var limit = 10;
            if (args.Length >= 2 && !int.TryParse(args[1], out limit))
            {
                Console.Error.WriteLine($"invalid limit: {args[1]}");
                return 1;
            }
            var events = store.Recent(limit);
            foreach (var ev in events) PrintEventSummary(ev);
            if (events.Count == 0) Console.WriteLine("(no events)");
            break;
        }
        case "count":
            g.Start();
            Console.WriteLine($"{store.Count()} events");
            break;
        case "verify":
        {
            g.Start();
            var result = store.VerifyChain();
            Console.WriteLine($"Chain verified: {result.Length} events, valid={result.Valid}");
            break;
        }
        case "head":
        {
            g.Start();
            var head = store.Head();
            if (head.IsNone)
                Console.WriteLine("(empty chain)");
            else
                PrintEvent(head.Unwrap());
            break;
        }
        case "help":
        case "-h":
        case "--help":
            Usage();
            break;
        default:
            Console.Error.WriteLine($"unknown command: {cmd}");
            Usage();
            return 1;
    }
}
finally
{
    g.Dispose();
}

return 0;

static void PrintEvent(Event ev)
{
    var obj = new Dictionary<string, object?>
    {
        ["id"] = ev.Id.Value,
        ["type"] = ev.Type.Value,
        ["source"] = ev.Source.Value,
        ["timestamp"] = DateTimeOffset.FromUnixTimeMilliseconds(ev.TimestampNanos / 1_000_000).ToString("o"),
        ["hash"] = ev.Hash.Value,
        ["prev_hash"] = ev.PrevHash.Value,
        ["conversation_id"] = ev.ConversationId.Value,
        ["version"] = ev.Version,
    };
    Console.WriteLine(JsonSerializer.Serialize(obj, new JsonSerializerOptions { WriteIndented = true }));
}

static void PrintEventSummary(Event ev)
{
    var ts = DateTimeOffset.FromUnixTimeMilliseconds(ev.TimestampNanos / 1_000_000).ToString("o");
    Console.WriteLine($"  {ev.Id.Value}  {ev.Type.Value}  {ts}");
}

static void Usage()
{
    Console.WriteLine("""
        eg — EventGraph CLI

        Usage:
          eg bootstrap              Initialize a new event graph
          eg get <event-id>         Get an event by ID
          eg recent [limit]         Show recent events (default: 10)
          eg count                  Show event count
          eg verify                 Verify hash chain integrity
          eg head                   Show the chain head
          eg help                   Show this help
        """);
}
