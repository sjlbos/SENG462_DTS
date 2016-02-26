
using System;
using System.Collections.Generic;
using System.Configuration;
using RabbitMessaging;
using ServiceHost;

namespace WorkloadGeneratorSlave
{
    public class WorkloadGeneratorService : WorkerHost
    {
        private int _numberOfHttpWorkers;
        private string _slaveName;

        protected override void InitializeService()
        {
            _numberOfHttpWorkers = Int32.Parse(ConfigurationManager.AppSettings["NumberOfHttpWorkers"]);
            _slaveName = ConfigurationManager.AppSettings["SlaveName"];
        }

        protected override IList<IWorker> GetWorkerList()
        {
            var receiver = RabbitMessengerFactory.GetReceiver("WorkloadQueueReceiver");
            //var publisher = RabbitMessengerFactory.GetPublisher("SlaveStatusPublisher");
            var worker = new WorkloadQueueMonitor(_slaveName, receiver, null, _numberOfHttpWorkers);
            return new List<IWorker>
            {
                worker
            };
        }
    }
}
