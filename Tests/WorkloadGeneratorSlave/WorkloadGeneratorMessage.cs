
namespace WorkloadGeneratorSlave
{
    internal abstract class WorkloadGeneratorMessage
    {
        public string Id { get; set; }
        public string MessageType { get; set; }
    }

    internal static class MessageType
    {
        public const string ControlMessage = "Control";
        public const string BatchOrderMessage = "BatchOrder";
    }
}
