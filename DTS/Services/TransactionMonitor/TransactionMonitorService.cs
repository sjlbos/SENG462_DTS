
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
        private int _queueReaderCount;
        private int _dbWriterCount;
        private string _dbConnectionString;
        private Uri _apiEndpoint;

        protected override void InitializeService()
        {
            _queueReaderCount = Int32.Parse(ConfigurationManager.AppSettings["NumberOfQueueReaders"]);
            _dbWriterCount = Int32.Parse(ConfigurationManager.AppSettings["NumberOfDbWriters"]);
            _apiEndpoint = new Uri(ConfigurationManager.AppSettings["ApiRoot"]);
            _dbConnectionString = ConfigurationManager.ConnectionStrings["DtsAuditDb"].ConnectionString;
        }

        protected override IList<IWorker> GetWorkerList()
        {
            var sharedBuffer = new BlockingCollection<TransactionEvent>();
            var workerList = new List<IWorker>();
            var repository = new PostgresAuditRepository(_dbConnectionString);

            for (int i = 0; i < _queueReaderCount; i++)
            {
                var receiver = RabbitMessengerFactory.GetReceiver("TransactionEventQueueReceiver");
                workerList.Add(new TransactionQueueMonitorWorker(String.Format(CultureInfo.InvariantCulture,
                    "Queue Monitor {0}", i), receiver, sharedBuffer));
            }
            for (int i = 0; i < _dbWriterCount; i++)
            {
                workerList.Add(new EventWriterWorker(String.Format(CultureInfo.InvariantCulture,
                    "Event Writer {0}", i), repository, sharedBuffer));
            }

            workerList.Add(new NancyHostLauncherWorker("Nancy Launcher", _apiEndpoint, repository));

            return workerList;
        }
    }
}
