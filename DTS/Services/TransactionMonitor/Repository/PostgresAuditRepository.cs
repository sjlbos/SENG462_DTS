
using System;
using System.Collections.Generic;
using System.Globalization;
using System.Net.Sockets;
using log4net;
using Npgsql;
using NpgsqlTypes;
using TransactionEvents;
using CommandType = System.Data.CommandType;

namespace TransactionMonitor.Repository
{
    public class PostgresAuditRepository : IAuditRepository
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (PostgresAuditRepository));

        private readonly string _connectionString;

        public PostgresAuditRepository(string connectionString)
        {
            if (connectionString == null)
                throw new ArgumentNullException("connectionString");

            _connectionString = connectionString;

            NpgsqlConnection.RegisterEnumGlobally<CommandType>("command_type");
            NpgsqlConnection.RegisterEnumGlobally<AccountAction>("account_action");
        }

        public void LogUserCommandEvent(UserCommandEvent userCommandEvent)
        {
            using (var command = new NpgsqlCommand("log_user_command_event"))
            {
                Log.DebugFormat("Inserting user command event {0} into databse...", userCommandEvent.Id);

                command.CommandType = CommandType.StoredProcedure;

                AddCommonEventPropertyParametersToCommand(command, userCommandEvent);

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Enum,
                    Value = userCommandEvent.CommandType
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Char,
                    Value = userCommandEvent.StockSymbol
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Money,
                    Value = userCommandEvent.Funds
                });

                int id = ExecuteInsertCommand(command);

                Log.DebugFormat(CultureInfo.InvariantCulture, 
                    "Successfully inserted user command event {0} (database id = {1}).", userCommandEvent.Id, id);
            }
        }

        public void LogQuoteServerEvent(QuoteServerEvent quoteServerEvent)
        {
            using (var command = new NpgsqlCommand("log_quote_server_event"))
            {
                Log.DebugFormat("Inserting quote server event {0} into databse...", quoteServerEvent.Id);

                command.CommandType = CommandType.StoredProcedure;

                AddCommonEventPropertyParametersToCommand(command, quoteServerEvent);

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Char,
                    Value = quoteServerEvent.StockSymbol
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Money,
                    Value = quoteServerEvent.Price
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.TimestampTZ,
                    Value = quoteServerEvent.QuoteServerTime
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Varchar,
                    Value = quoteServerEvent.CryptoKey
                });

                int id = ExecuteInsertCommand(command);

                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Successfully inserted quote server event {0} (database id = {1}).", quoteServerEvent.Id, id);
            }
        }

        public void LogAccountTransactionEvent(AccountTransactionEvent accountTransactionEvent)
        {
            using (var command = new NpgsqlCommand("log_account_transaction_event"))
            {
                Log.DebugFormat("Inserting account transaction event {0} into databse...", accountTransactionEvent.Id);

                command.CommandType = CommandType.StoredProcedure;

                AddCommonEventPropertyParametersToCommand(command, accountTransactionEvent);

                command.Parameters.Add(new NpgsqlParameter
                {
                    Value = accountTransactionEvent.AccountAction
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Money,
                    Value = accountTransactionEvent.Funds
                });

                int id = ExecuteInsertCommand(command);

                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Successfully inserted account transaction event {0} (database id = {1}).", accountTransactionEvent.Id, id);
            }
        }

        public void LogSystemEvent(SystemEvent systemEvent)
        {
            using (var command = new NpgsqlCommand("log_system_event"))
            {
                Log.DebugFormat("Inserting system event {0} into databse...", systemEvent.Id);

                command.CommandType = CommandType.StoredProcedure;

                AddCommonEventPropertyParametersToCommand(command, systemEvent);

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Enum,
                    Value = systemEvent.Command
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Char,
                    Value = systemEvent.StockSymbol
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Money,
                    Value = systemEvent.Funds
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Varchar,
                    Value = systemEvent.FileName
                });

                int id = ExecuteInsertCommand(command);

                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Successfully inserted system event {0} (database id = {1}).", systemEvent.Id, id);
            }
        }

        public void LogErrorEvent(ErrorEvent errorEvent)
        {
            using (var command = new NpgsqlCommand("log_error_event"))
            {
                Log.DebugFormat("Inserting error event {0} into databse...", errorEvent.Id);

                command.CommandType = CommandType.StoredProcedure;

                AddCommonEventPropertyParametersToCommand(command, errorEvent);

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Enum,
                    Value = errorEvent.Command
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Char,
                    Value = errorEvent.StockSymbol
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Money,
                    Value = errorEvent.Funds
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Varchar,
                    Value = errorEvent.ErrorMessage
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Varchar,
                    Value = errorEvent.FileName
                });

                int id = ExecuteInsertCommand(command);

                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Successfully inserted error event {0} (database id = {1}).", errorEvent.Id, id);
            }
        }

        public void LogDebugEvent(DebugEvent debugEvent)
        {
            using (var command = new NpgsqlCommand("log_debug_event"))
            {
                Log.DebugFormat("Inserting debug event {0} into databse...", debugEvent.Id);

                command.CommandType = CommandType.StoredProcedure;

                AddCommonEventPropertyParametersToCommand(command, debugEvent);

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Enum,
                    Value = debugEvent.Command
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Char,
                    Value = debugEvent.StockSymbol
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Money,
                    Value = debugEvent.Funds
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Varchar,
                    Value = debugEvent.FileName
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Varchar,
                    Value = debugEvent.DebugMessage
                });

                int id = ExecuteInsertCommand(command);

                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Successfully inserted debug event {0} (database id = {1}).", debugEvent.Id, id);
            }
        }

        public IEnumerable<TransactionEvent> GetLogOfTransaction(Guid transactionId)
        {
            throw new NotImplementedException();
        }

        public IEnumerable<TransactionEvent> GetAllLogsBeween(DateTime start, DateTime end)
        {
            throw new NotImplementedException();
        }

        public IEnumerable<TransactionEvent> GetAllLogs()
        {
            throw new NotImplementedException();
        }

        private void AddCommonEventPropertyParametersToCommand(NpgsqlCommand command, TransactionEvent transactionEvent)
        {
            command.Parameters.Add(new NpgsqlParameter
            {
                NpgsqlDbType = NpgsqlDbType.TimestampTZ,
                Value = transactionEvent.OccuredAt
            });

            command.Parameters.Add(new NpgsqlParameter
            {
                NpgsqlDbType = NpgsqlDbType.Uuid,
                Value = transactionEvent.TransactionId
            });

            command.Parameters.Add(new NpgsqlParameter
            {
                NpgsqlDbType = NpgsqlDbType.Char,
                Value = transactionEvent.UserId
            });

            command.Parameters.Add(new NpgsqlParameter
            {
                NpgsqlDbType = NpgsqlDbType.Varchar,
                Value = transactionEvent.Service
            });

            command.Parameters.Add(new NpgsqlParameter
            {
                NpgsqlDbType = NpgsqlDbType.Varchar,
                Value = transactionEvent.Server
            });
        }

        private int ExecuteInsertCommand(NpgsqlCommand command)
        {
            try
            {
                using (var connection = new NpgsqlConnection(_connectionString))
                {
                    connection.Open();
                    command.Connection = connection;
                    return (int) command.ExecuteScalar();
                }
            }
            catch (SocketException ex)
            {
                throw new RepositoryException(String.Format(CultureInfo.InvariantCulture,
                    "Encountered an error while attempting to execute database command \"{0}\" using connection string \"{1}\".", command.CommandText, _connectionString), ex);
            }
            catch (NpgsqlException ex)
            {
                throw new RepositoryException(String.Format(CultureInfo.InvariantCulture,
                    "Encountered an error while attempting to execute database command \"{0}\" using connection string \"{1}\".", command.CommandText, _connectionString), ex);
            }
        }
    }
}
