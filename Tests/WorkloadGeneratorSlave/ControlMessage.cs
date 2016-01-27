
namespace WorkloadGeneratorSlave
{
    internal class ControlMessage : WorkloadGeneratorMessage
    {
        public const string StartCommand = "Start";

        public string Command { get; set; }
    }
}
