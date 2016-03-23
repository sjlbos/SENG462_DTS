
using System;
using System.Collections.Generic;
using log4net;
using MongoDB.Driver;
using TransactionEvents;

namespace TransactionMonitor.Repository
{
    public class MongoDbAuditRepository : IAuditRepository
    {
        private static ILog Log = LogManager.GetLogger(typeof (MongoDbAuditRepository));

        private readonly IMongoClient _mongoClient;
        private readonly IMongoDatabase _auditDb;

        public MongoDbAuditRepository(string connectionString)
        {
            if(connectionString == null)
                throw new ArgumentNullException("connectionString");

            string databaseName = MongoUrl.Create(connectionString).DatabaseName;
            _mongoClient = new MongoClient(connectionString);
            _auditDb = _mongoClient.GetDatabase(databaseName);
        }

        public void LogUserCommandEvent(UserCommandEvent userCommandEvent)
        {
            throw new NotImplementedException();
        }

        public void LogQuoteServerEvent(QuoteServerEvent quoteServerEvent)
        {
            throw new NotImplementedException();
        }

        public void LogAccountTransactionEvent(AccountTransactionEvent accountTransactionEvent)
        {
            throw new NotImplementedException();
        }

        public void LogSystemEvent(SystemEvent systemEvent)
        {
            throw new NotImplementedException();
        }

        public void LogErrorEvent(ErrorEvent errorEvent)
        {
            throw new NotImplementedException();
        }

        public void LogDebugEvent(DebugEvent debugEvent)
        {
            throw new NotImplementedException();
        }

        public IEnumerable<TransactionEvent> GetLogsForUser(string userId, DateTime start, DateTime end)
        {
            throw new NotImplementedException();
        }

        public IEnumerable<TransactionEvent> GetAllLogs(DateTime start, DateTime end)
        {
            throw new NotImplementedException();
        }
    }
}
