
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
            int totalWorkerCount = _dbWriterCount + _queueReaderCount;
            var workerList = new List<IWorker>(totalWorkerCount);
            var repository = new PostgresAuditRepository(_dbConnectionString);

            for (int i = 0; i < _queueReaderCount; i++)
            {
                var receiver = RabbitMessengerFactory.GetReceiver("TransactionEventQueueReceiver");
                workerList.Add(new TransactionQueueMonitorWorker(
                    i.ToString(CultureInfo.InvariantCulture), receiver, sharedBuffer));
            }
            for (int i = _queueReaderCount; i < totalWorkerCount; i++)
            {
                workerList.Add(new EventWriterWorker(
                    i.ToString(CultureInfo.InvariantCulture), repository, sharedBuffer));
            }

            workerList.Add(new NancyHostLauncherWorker("Nancy Launcher", _apiEndpoint, repository));

            return workerList;
        }
    }
}
