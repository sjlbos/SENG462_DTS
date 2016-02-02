
using System;
using System.Collections.Generic;
using System.Configuration;
using RabbitMessaging;
using ServiceHost;

namespace WorkloadGeneratorSlave
{
    public class WorkloadGeneratorService : WorkerHost
    {
        private int _numberOfWorkers;

        protected override void InitializeService()
        {
            _numberOfWorkers = Int32.Parse(ConfigurationManager.AppSettings["NumberOfWorkers"]);
        }

        protected override IList<IWorker> GetWorkerList()
        {
            var receiver = RabbitMessengerFactory.GetReceiver("WorkloadQueueReceiver");
            var worker = new WorkloadQueueMonitor("1", receiver, _numberOfWorkers);
            return new[] {worker};
        }
    }
}
