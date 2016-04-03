
using System;
using System.Collections.Generic;
using System.Configuration;
using System.Globalization;
using RabbitMessaging;
using ServiceHost;
using TransactionMonitor.Api;
using TransactionMonitor.Repository;

namespace TransactionMonitor
{
    public class TransactionMonitorService : WorkerHost
    {
        private int _workerCount;
        private Uri _apiEndpoint;

        private IAuditRepository _repository;

        protected override void InitializeService()
        {
            _workerCount = Int32.Parse(ConfigurationManager.AppSettings["WorkerCount"]);
            _apiEndpoint = new Uri(ConfigurationManager.AppSettings["ApiRoot"]);

            string databaseToUse = ConfigurationManager.AppSettings["UseDatabase"];

            if (String.Equals("MongoDB", databaseToUse))
            {
                string connnectionString = ConfigurationManager.ConnectionStrings["MongoAuditDB"].ConnectionString;
                string collectionName = ConfigurationManager.AppSettings["EventCollectionName"];
                _repository = new MongoDbAuditRepository(connnectionString, collectionName);
            }
            else if (String.Equals("Postgres", databaseToUse))
            {
                string connectionString = ConfigurationManager.ConnectionStrings["PostgresAuditDB"].ConnectionString;
                _repository = new PostgresAuditRepository(connectionString);
            }
            else
            {
                throw new ConfigurationErrorsException("Invalid database type.");
            }
        }

        protected override IList<IWorker> GetWorkerList()
        {
            var workerList = new List<IWorker>();

            for (int i = 0; i < _workerCount; i++)
            {
                var receiver = RabbitMessengerFactory.GetReceiver("TransactionEventQueueReceiver");
                workerList.Add(new TransactionMonitorWorker(String.Format(CultureInfo.InvariantCulture,
                    "Worker {0}", i), receiver, _repository));
            }

            workerList.Add(new NancyHostLauncherWorker("Nancy Launcher", _apiEndpoint, _repository));

            return workerList;
        }
    }
}
