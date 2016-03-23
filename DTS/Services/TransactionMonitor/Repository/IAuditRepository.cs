using System;
using System.Collections.Generic;
using TransactionEvents;

namespace TransactionMonitor.Repository
{
    public interface IAuditRepository
    {
        void LogTransactionEvent(TransactionEvent transactionEvent);
        IEnumerable<TransactionEvent> GetLogsForUser(string userId, DateTime start, DateTime end);
        IEnumerable<TransactionEvent> GetAllLogs(DateTime start, DateTime end);
    }
}
