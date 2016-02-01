
using System.Collections.Generic;
using System.Configuration;
using RabbitMessaging;
using ServiceHost;

namespace WorkloadGeneratorSlave
{
    public class WorkloadGeneratorService : WorkerHost<WorkloadQueueMonitor>
    {
        private int _numberOfWorkers;

        protected override void InitializeService()
        {
            _numberOfWorkers = int.Parse(ConfigurationManager.AppSettings["NumberOfWorkers"]);
        }

        protected override IList<WorkloadQueueMonitor> GetWorkerList()
        {
            var receiver = RabbitMessengerFactory.GetReceiver("WorkloadQueueReceiver");
            var worker = new WorkloadQueueMonitor("1", receiver, _numberOfWorkers);
            return new[] {worker};
        }
    }
}
