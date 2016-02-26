
using System;

namespace WorkloadGeneratorSlave
{
    public class StatusMessage
    {
        public string SlaveName { get; set; }
        public string Status { get; set; }
        public DateTime Timestamp { get; set; }
    }
}
