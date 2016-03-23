
using System;
using System.Collections.Generic;
using System.Globalization;
using System.Net.Sockets;
using log4net;
using Npgsql;
using NpgsqlTypes;
using TransactionEvents;

namespace TransactionMonitor.Repository
{
    public class PostgresAuditRepository : IAuditRepository
    {
        private const string UserCommandEvent = "UserCommandEvent";
        private const string QuoteServerEvent = "QuoteServerEvent";
        private const string AccountTransactionEvent = "AccountTransactionEvent";
        private const string SystemEvent = "SystemEvent";
        private const string ErrorEvent = "ErrorEvent";
        private const string DebugEvent = "DebugEvent";

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

        public void LogTransactionEvent(TransactionEvent transactionEvent)
        {
            if (transactionEvent is UserCommandEvent)
            {
                LogUserCommandEvent(transactionEvent as UserCommandEvent);
                return;
            }

            if (transactionEvent is QuoteServerEvent)
            {
                LogQuoteServerEvent(transactionEvent as QuoteServerEvent);
                return;
            }

            if (transactionEvent is AccountTransactionEvent)
            {
                LogAccountTransactionEvent(transactionEvent as AccountTransactionEvent);
                return;
            }

            if (transactionEvent is SystemEvent)
            {
                LogSystemEvent(transactionEvent as SystemEvent);
                return;
            }

            if (transactionEvent is ErrorEvent)
            {
                LogErrorEvent(transactionEvent as ErrorEvent);
                return;
            }

            if (transactionEvent is DebugEvent)
            {
                LogDebugEvent(transactionEvent as DebugEvent);
                return;
            }
        }

        public void LogUserCommandEvent(UserCommandEvent userCommandEvent)
        {
            using (var command = new NpgsqlCommand("log_user_command_event"))
            {
                Log.DebugFormat("Inserting user command event {0} into databse...", userCommandEvent.Id);

                command.CommandType = System.Data.CommandType.StoredProcedure;

                AddCommonEventPropertyParametersToCommand(command, userCommandEvent);

                command.Parameters.Add(new NpgsqlParameter
                {
                    Value = userCommandEvent.Command
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Char,
                    Value =  ((object) userCommandEvent.StockSymbol) ?? DBNull.Value
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Money,
                    Value = ((object) userCommandEvent.Funds) ?? DBNull.Value
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

                command.CommandType = System.Data.CommandType.StoredProcedure;

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

                command.CommandType = System.Data.CommandType.StoredProcedure;

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

                command.CommandType = System.Data.CommandType.StoredProcedure;

                AddCommonEventPropertyParametersToCommand(command, systemEvent);

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Enum,
                    Value = systemEvent.Command
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Char,
                    Value = ((object) systemEvent.StockSymbol) ?? DBNull.Value
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Money,
                    Value = ((object) systemEvent.Funds) ?? DBNull.Value
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Varchar,
                    Value = ((object) systemEvent.FileName) ?? DBNull.Value
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

                command.CommandType = System.Data.CommandType.StoredProcedure;

                AddCommonEventPropertyParametersToCommand(command, errorEvent);

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Enum,
                    Value = errorEvent.Command
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Char,
                    Value = ((object) errorEvent.StockSymbol) ?? DBNull.Value
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Money,
                    Value = ((object) errorEvent.Funds) ?? DBNull.Value
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Varchar,
                    Value = ((object) errorEvent.ErrorMessage) ?? DBNull.Value
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Varchar,
                    Value = ((object) errorEvent.FileName) ?? DBNull.Value
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

                command.CommandType = System.Data.CommandType.StoredProcedure;

                AddCommonEventPropertyParametersToCommand(command, debugEvent);

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Enum,
                    Value = debugEvent.Command
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Char,
                    Value = ((object) debugEvent.StockSymbol) ?? DBNull.Value
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Money,
                    Value = ((object) debugEvent.Funds) ?? DBNull.Value
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Varchar,
                    Value = ((object) debugEvent.FileName) ?? DBNull.Value
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Varchar,
                    Value = ((object) debugEvent.DebugMessage) ?? DBNull.Value
                });

                int id = ExecuteInsertCommand(command);

                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Successfully inserted debug event {0} (database id = {1}).", debugEvent.Id, id);
            }
        }

        public IEnumerable<TransactionEvent> GetLogsForUser(string userId, DateTime start, DateTime end)
        {
            if(userId == null)
                throw new ArgumentNullException("userId");

            var command = new NpgsqlCommand("get_all_events_by_user");
            command.CommandType = System.Data.CommandType.StoredProcedure;

            command.Parameters.Add(new NpgsqlParameter
            {
                NpgsqlDbType = NpgsqlDbType.Char,
                Value = userId
            });

            command.Parameters.Add(new NpgsqlParameter
            {
                NpgsqlDbType = NpgsqlDbType.TimestampTZ,
                Value = start
            });

            command.Parameters.Add(new NpgsqlParameter
            {
                NpgsqlDbType = NpgsqlDbType.TimestampTZ,
                Value = end
            });

            Log.DebugFormat(CultureInfo.InvariantCulture,
                "Querying all event logs for user \"{0}\" using query \"{1}\".", userId, command.CommandText);

            var results = GetTransactionEventsUsingQuery(command);

            Log.DebugFormat(CultureInfo.InvariantCulture,
                "Query \"{0}\" completed successfully.", command.CommandText);

            return results;
        }

        public IEnumerable<TransactionEvent> GetAllLogs(DateTime start, DateTime end)
        {
            var command = new NpgsqlCommand("get_all_events");
                command.CommandType = System.Data.CommandType.StoredProcedure;

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.TimestampTZ,
                    Value = start
                });

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.TimestampTZ,
                    Value = end
                });

                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Querying all event logs using query \"{0}\".", command.CommandText);

                var results = GetTransactionEventsUsingQuery(command);

                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Query \"{0}\" completed successfully.", command.CommandText);

                return results;
        }

        private static void AddCommonEventPropertyParametersToCommand(NpgsqlCommand command, TransactionEvent transactionEvent)
        {
            command.Parameters.Add(new NpgsqlParameter
            {
                NpgsqlDbType = NpgsqlDbType.TimestampTZ,
                Value = transactionEvent.OccuredAt
            });

            command.Parameters.Add(new NpgsqlParameter
            {
                NpgsqlDbType = NpgsqlDbType.Integer,
                Value = transactionEvent.TransactionId
            });

            command.Parameters.Add(new NpgsqlParameter
            {
                NpgsqlDbType = NpgsqlDbType.Char,
                Value = ((object)transactionEvent.UserId) ?? DBNull.Value
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

        private IEnumerable<TransactionEvent> GetTransactionEventsUsingQuery(NpgsqlCommand command)
        {
            NpgsqlConnection connection = null;
            NpgsqlDataReader reader = null;

            try
            {
                connection = new NpgsqlConnection(_connectionString);
                connection.Open();
                command.Connection = connection;
                reader = command.ExecuteReader();
                if (reader.HasRows)
                {
                    while (reader.Read())
                    {
                        yield return GetTransactionEventFromRecord(reader);
                    }
                }
            }
            finally
            {
                if (reader != null)
                {
                    reader.Close();
                    reader.Dispose();
                }

                if (command != null)
                {
                    command.Dispose();
                }

                if (connection != null)
                {
                    connection.Close();
                    connection.Dispose();
                }
            } 
        }

        /* Expected record:
         *      (0)  id
         *      (1)  logged_at
         *      (2)  occured_at
         *      (3)  type
         *      (4)  transaction_id
         *      (5)  user_id
         *      (6)  service
         *      (7)  server
         *      (8)  command
         *      (9)  stock
         *      (10) funds
         *      (11) filename
         *      (12) message
         *      (13) action
         *      (14) quote_server_time
         *      (15) cryptokey
         */
        private TransactionEvent GetTransactionEventFromRecord(NpgsqlDataReader reader)
        {
            string eventType = reader.GetString(3);
            switch (eventType)
            {
                case UserCommandEvent:
                    return GetUserCommandEventFromRecord(reader);
                case QuoteServerEvent:
                    return GetQuoteServerEventFromRecord(reader);
                case AccountTransactionEvent:
                    return GetAccountTransactionEventFromRecord(reader);
                case SystemEvent:
                    return GetSystemEventFromRecord(reader);
                case ErrorEvent:
                    return GetErrorEventFromRecord(reader);
                case DebugEvent:
                    return GetDebugEventFromRecord(reader);
                default:
                    throw new UnrecognizedTransactionEventException("Received an unrecognized event type from the database: " + eventType);
            }
        }

        private void FillBaseEventPropertiesFromRecord(TransactionEvent te, NpgsqlDataReader reader)
        {
            te.LoggedAt = reader.GetDateTime(1);
            te.OccuredAt = reader.GetDateTime(2);
            te.TransactionId = reader.GetInt32(4);
            te.UserId = (reader.IsDBNull(5)) ? null : reader.GetString(5);
            te.Service = reader.GetString(6);
            te.Server = reader.GetString(7);
        }

        private UserCommandEvent GetUserCommandEventFromRecord(NpgsqlDataReader reader)
        {
            var userCommandEvent = new UserCommandEvent
            {
                Command = (CommandType) reader.GetValue(8),
                StockSymbol = (reader.IsDBNull(9)) ? null : reader.GetString(9),
                Funds = (reader.IsDBNull(10)) ? null : (decimal?) reader.GetDecimal(10)
            };
            FillBaseEventPropertiesFromRecord(userCommandEvent, reader);
            return userCommandEvent;
        }

        private QuoteServerEvent GetQuoteServerEventFromRecord(NpgsqlDataReader reader)
        {
            var quoteServerEvent = new QuoteServerEvent
            {
                StockSymbol = reader.GetString(9),
                Price = reader.GetDecimal(10),
                QuoteServerTime = reader.GetDateTime(14),
                CryptoKey = reader.GetString(15)
            };
            FillBaseEventPropertiesFromRecord(quoteServerEvent, reader);
            return quoteServerEvent;
        }

        private AccountTransactionEvent GetAccountTransactionEventFromRecord(NpgsqlDataReader reader)
        {
            var accountTransactionEvent = new AccountTransactionEvent
            {
                AccountAction = (AccountAction) reader.GetValue(13),
                Funds = reader.GetDecimal(10)
            };
            FillBaseEventPropertiesFromRecord(accountTransactionEvent, reader);
            return accountTransactionEvent;
        }

        private SystemEvent GetSystemEventFromRecord(NpgsqlDataReader reader)
        {
            var systemEvent = new SystemEvent
            {
                Command = (CommandType)reader.GetValue(8),
                StockSymbol = (reader.IsDBNull(9)) ? null : reader.GetString(9),
                Funds = (reader.IsDBNull(10)) ? null : (decimal?) reader.GetDecimal(10),
                FileName = (reader.IsDBNull(11)) ? null : reader.GetString(11)
            };
            FillBaseEventPropertiesFromRecord(systemEvent, reader);
            return systemEvent;
        }

        private ErrorEvent GetErrorEventFromRecord(NpgsqlDataReader reader)
        {
            var errorEvent = new ErrorEvent
            {
                Command = (CommandType)reader.GetValue(8),
                StockSymbol = (reader.IsDBNull(9)) ? null : reader.GetString(9),
                Funds = (reader.IsDBNull(10)) ? null : (decimal?) reader.GetDecimal(10),
                FileName = (reader.IsDBNull(11)) ? null : reader.GetString(11),
                ErrorMessage = (reader.IsDBNull(12)) ? null : reader.GetString(12)
            };
            FillBaseEventPropertiesFromRecord(errorEvent, reader);
            return errorEvent;
        }

        private DebugEvent GetDebugEventFromRecord(NpgsqlDataReader reader)
        {
            var debugEvent = new DebugEvent
            {
                Command = (CommandType)reader.GetValue(8),
                StockSymbol = (reader.IsDBNull(9)) ? null : reader.GetString(9),
                Funds = (reader.IsDBNull(10)) ? null : (decimal?) reader.GetDecimal(10),
                FileName = (reader.IsDBNull(11)) ? null : reader.GetString(11),
                DebugMessage = (reader.IsDBNull(12)) ? null : reader.GetString(12)
            };
            FillBaseEventPropertiesFromRecord(debugEvent, reader);
            return debugEvent;
        }
    }
}
