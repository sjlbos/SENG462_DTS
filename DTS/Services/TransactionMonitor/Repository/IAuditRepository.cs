using System;
using System.Collections.Generic;
using TransactionEvents;

namespace TransactionMonitor.Repository
{
    public interface IAuditRepository
    {
        void LogUserCommandEvent(UserCommandEvent userCommandEvent);
        void LogQuoteServerEvent(QuoteServerEvent quoteServerEvent);
        void LogAccountTransactionEvent(AccountTransactionEvent accountTransactionEvent);
        void LogSystemEvent(SystemEvent systemEvent);
        void LogErrorEvent(ErrorEvent errorEvent);
        void LogDebugEvent(DebugEvent debugEvent);

        IEnumerable<TransactionEvent> GetLogOfTransaction(Guid transactionId);
        IEnumerable<TransactionEvent> GetAllLogsBeween(DateTime start, DateTime end);
        IEnumerable<TransactionEvent> GetAllLogs();
    }
}
