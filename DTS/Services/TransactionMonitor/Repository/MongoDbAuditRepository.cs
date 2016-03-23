
using System;
using System.Collections.Generic;
using log4net;
using MongoDB.Bson.Serialization;
using MongoDB.Driver;
using TransactionEvents;

namespace TransactionMonitor.Repository
{
    public class MongoDbAuditRepository : IAuditRepository
    {
        private static ILog Log = LogManager.GetLogger(typeof (MongoDbAuditRepository));

        private readonly IMongoCollection<TransactionEvent> _eventCollection; 

        public MongoDbAuditRepository(string connectionString, string eventCollectionName)
        {
            if(connectionString == null)
                throw new ArgumentNullException("connectionString");

            RegisterEventClasses();

            string databaseName = MongoUrl.Create(connectionString).DatabaseName;
            var mongoClient = new MongoClient(connectionString);
            var auditDb = mongoClient.GetDatabase(databaseName);
            _eventCollection = auditDb.GetCollection<TransactionEvent>(eventCollectionName);
        }

        private static void RegisterEventClasses()
        {
            BsonClassMap.RegisterClassMap<TransactionEvent>(cm =>
            {
                cm.AutoMap();
                cm.SetIsRootClass(true);
            });
            BsonClassMap.RegisterClassMap<UserCommandEvent>();
            BsonClassMap.RegisterClassMap<QuoteServerEvent>();
            BsonClassMap.RegisterClassMap<AccountTransactionEvent>();
            BsonClassMap.RegisterClassMap<SystemEvent>();
            BsonClassMap.RegisterClassMap<ErrorEvent>();
            BsonClassMap.RegisterClassMap<DebugEvent>();
        }

        #region IAuditRepository

        public void LogUserCommandEvent(UserCommandEvent userCommandEvent)
        {
            _eventCollection.InsertOne(userCommandEvent);
        }

        public void LogQuoteServerEvent(QuoteServerEvent quoteServerEvent)
        {
            _eventCollection.InsertOne(quoteServerEvent);
        }

        public void LogAccountTransactionEvent(AccountTransactionEvent accountTransactionEvent)
        {
            _eventCollection.InsertOne(accountTransactionEvent);
        }

        public void LogSystemEvent(SystemEvent systemEvent)
        {
            _eventCollection.InsertOne(systemEvent);
        }

        public void LogErrorEvent(ErrorEvent errorEvent)
        {
            _eventCollection.InsertOne(errorEvent);
        }

        public void LogDebugEvent(DebugEvent debugEvent)
        {
            _eventCollection.InsertOne(debugEvent);
        }

        public IEnumerable<TransactionEvent> GetLogsForUser(string userId, DateTime start, DateTime end)
        {
            var builder = Builders<TransactionEvent>.Filter;
            var filter = builder.Eq(e => e.UserId, userId) & builder.Gt(e => e.OccuredAt, start) & builder.Lt(e => e.OccuredAt, end);
            return _eventCollection.Find(filter).ToEnumerable();
        }

        public IEnumerable<TransactionEvent> GetAllLogs(DateTime start, DateTime end)
        {
            var builder = Builders<TransactionEvent>.Filter;
            var filter = builder.Gt(e => e.OccuredAt, start) & builder.Lt(e => e.OccuredAt, end);
            return _eventCollection.Find(filter).ToEnumerable();
        }

        #endregion
    }
}
