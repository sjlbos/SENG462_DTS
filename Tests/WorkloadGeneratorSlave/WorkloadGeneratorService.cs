
using System.Collections.Generic;
using System.Configuration;
using System.Globalization;
using RabbitMessaging;
using ServiceHost;

namespace WorkloadGeneratorSlave
{
    public class WorkloadGeneratorService : WorkerHost<WorkloadGeneratorWorker>
    {
        private int _numberOfWorkers;

        protected override void InitializeService()
        {
            _numberOfWorkers = int.Parse(ConfigurationManager.AppSettings["NumberOfWorkers"]);
        }

        protected override IList<WorkloadGeneratorWorker> GetWorkerList()
        {
            var workers = new List<WorkloadGeneratorWorker>(_numberOfWorkers);
            for (int i = 0; i < _numberOfWorkers; i++)
            {
                var receiver = RabbitMessengerFactory.GetReceiver("WorkloadQueueReceiver");
                var worker = new WorkloadGeneratorWorker(i.ToString(CultureInfo.InvariantCulture), receiver);
                workers.Add(worker);   
            }
            return workers;
        }
    }
}
