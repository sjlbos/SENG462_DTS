
using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Configuration;
using System.Globalization;
using RabbitMessaging;
using ServiceHost;
using TransactionEvents;
using TransactionMonitor.Api;
using TransactionMonitor.Repository;

namespace TransactionMonitor
{
    public class TransactionMonitorService : WorkerHost
    {
        private int _workerCount;
        private string _dbConnectionString;
        private Uri _apiEndpoint;

        protected override void InitializeService()
        {
            _workerCount = Int32.Parse(ConfigurationManager.AppSettings["WorkerCount"]);
            _apiEndpoint = new Uri(ConfigurationManager.AppSettings["ApiRoot"]);
            _dbConnectionString = ConfigurationManager.ConnectionStrings["DtsAuditDb"].ConnectionString;
        }

        protected override IList<IWorker> GetWorkerList()
        {
     
            var workerList = new List<IWorker>();
            var repository = new PostgresAuditRepository(_dbConnectionString);

            for (int i = 0; i < _workerCount; i++)
            {
                var receiver = RabbitMessengerFactory.GetReceiver("TransactionEventQueueReceiver");
                workerList.Add(new TransactionMonitorWorker(String.Format(CultureInfo.InvariantCulture,
                    "Worker {0}", i), receiver, repository));
            }

            workerList.Add(new NancyHostLauncherWorker("Nancy Launcher", _apiEndpoint, repository));

            return workerList;
        }
    }
}
