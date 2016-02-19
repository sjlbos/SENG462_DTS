
using System;
using System.Collections.Generic;
using System.Data;
using System.Globalization;
using System.Net.Sockets;
using log4net;
using Npgsql;
using NpgsqlTypes;
using TriggerManager.Models;

namespace TriggerManager.Repository
{
    /// <summary>
    /// A class to provide a read-only abstraction layer around the DTS database's trigger data.
    /// </summary>
    public class PostgresTriggerRepository : ITriggerRepository
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (PostgresTriggerRepository));

        private readonly string _connectionString;

        public PostgresTriggerRepository(string connectionString)
        {
            if (connectionString == null)
                throw new ArgumentNullException("connectionString");

            _connectionString = connectionString;

            NpgsqlConnection.RegisterEnumGlobally<TriggerType>("trigger_type");
        }

        #region ITriggerRepository

        /// <summary>
        /// Returns all sell triggers belonging to the user with the specified user ID.
        /// </summary>
        /// <param name="userId">The user ID of a DTS user.</param>
        /// <returns>A list of Trigger objects or an empty list if no triggers exist for the specified user.</returns>
        public IList<Trigger> GetBuyTriggersForUser(int userId)
        {
            using (var command = new NpgsqlCommand("get_user_buy_triggers"))
            {
                command.CommandType = CommandType.StoredProcedure;

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Integer,
                    Value = userId
                });

                return ExecuteGetTriggerCommand(command);
            }
        }

        /// <summary>
        /// Returns all sell triggers belonging to the user with the specified user ID.
        /// </summary>
        /// <param name="userId">The user ID of a DTS user.</param>
        /// <returns>A list of Trigger objects or an empty list if no triggers exist for the specfied user.</returns>
        public IList<Trigger> GetSellTriggersForUser(int userId)
        {
            using (var command = new NpgsqlCommand("get_user_sell_triggers"))
            {
                command.CommandType = CommandType.StoredProcedure;

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Integer,
                    Value = userId
                });

                return ExecuteGetTriggerCommand(command);
            }
        }

        /// <summary>
        /// Returns a list containing all buy and sell triggers belonging to all DTS users.
        /// </summary>
        /// <returns>A list of Trigger objects or an empty list if no triggers exist in the DTS.</returns>
        public IList<Trigger> GetAllTriggers()
        {
            using (var command = new NpgsqlCommand("get_all_triggers"))
            {
                command.CommandType = CommandType.StoredProcedure;

                return ExecuteGetTriggerCommand(command);
            }
        }

        #endregion

        #region Helper Methods

        private IList<Trigger> ExecuteGetTriggerCommand(NpgsqlCommand command)
        {
            try
            {
                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Executing database query \"{0}\" using connection string \"{1}\"...", command.CommandText, _connectionString);

                using (var connection = new NpgsqlConnection(_connectionString))
                {
                    connection.Open();
                    command.Connection = connection;
                    using (var reader = command.ExecuteReader())
                    {
                        Log.Debug("Query executed successfully.");

                        var triggerList = new List<Trigger>();
                        if (reader.HasRows)
                        {
                            while (reader.Read())
                            {
                                triggerList.Add(GetTriggerFromRecord(reader));
                            }
                        }
                        return triggerList;
                    }
                }
            }
            catch (SocketException ex)
            {
                throw new RepositoryException(String.Format(CultureInfo.InvariantCulture,
                    "Encountered an error while attempting to execute database command \"{0}\" using connection string \"{1}\".", 
                    command.CommandText, _connectionString), ex);
            }
            catch (NpgsqlException ex)
            {
                throw new RepositoryException(String.Format(CultureInfo.InvariantCulture,
                    "Encountered an error while attempting to execute database command \"{0}\" using connection string \"{1}\".",
                    command.CommandText, _connectionString), ex);
            }
        }

        /* Expected record:
         *      (0) id 
         *      (1) uid
         *      (2) stock
         *      (3) type
         *      (4) trigger_price
         *      (5) num_shares
         *      (6) created_at
         */
        private static Trigger GetTriggerFromRecord(IDataRecord reader)
        {
            return new Trigger
            {
                Id = reader.GetInt32(0),
                UserId = reader.GetInt32(1),
                StockSymbol = reader.GetString(2),
                Type = (TriggerType) reader.GetValue(3),
                TriggerPrice = reader.GetDecimal(4),
                NumberOfShares = reader.GetInt32(5)
            };
        }

        #endregion
    }
}
