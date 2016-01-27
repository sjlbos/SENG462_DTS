
using System.Collections.Generic;

namespace WorkloadGeneratorSlave
{
    internal class WorkloadBatchMessage : WorkloadGeneratorMessage
    {
        public IList<ApiCommand> Commands { get; set; } 
    }
}
