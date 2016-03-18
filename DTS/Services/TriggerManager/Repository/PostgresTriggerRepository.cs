
using System;
using System.Collections.Generic;
using System.Data;
using System.Globalization;
using System.Linq;
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

        private readonly IList<string> _connectionStrings;

        public PostgresTriggerRepository(IList<string> connectionStrings)
        {
            if (connectionStrings == null)
                throw new ArgumentNullException("connectionStrings");

            _connectionStrings = connectionStrings;

            NpgsqlConnection.RegisterEnumGlobally<TriggerType>("trigger_type");
        }

        #region ITriggerRepository

        /// <summary>
        /// Returns all sell triggers belonging to the user with the specified user ID.
        /// </summary>
        /// <param name="userDbId">The database ID of a DTS user.</param>
        /// <param name="userId">The user ID string of a DTS user.</param>
        /// <returns>A list of Trigger objects or an empty list if no triggers exist for the specified user.</returns>
        public IList<Trigger> GetBuyTriggersForUser(int userDbId, string userId)
        {
            if (userId == null)
                throw new ArgumentNullException("userId");

            using (var command = new NpgsqlCommand("get_user_buy_triggers"))
            {
                command.CommandType = CommandType.StoredProcedure;

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Integer,
                    Value = userId
                });

                return ExecuteGetTriggerCommand(command, GetDatabaseInstanceFromUserId(userId), userId);
            }
        }

        /// <summary>
        /// Returns all sell triggers belonging to the user with the specified user ID.
        /// </summary>
        /// <param name="userDbId">The database ID of a DTS user.</param>
        /// <param name="userId">The user ID string of a DTS user.</param>
        /// <returns>A list of Trigger objects or an empty list if no triggers exist for the specfied user.</returns>
        public IList<Trigger> GetSellTriggersForUser(int userDbId, string userId)
        {
            if (userId == null)
                throw new ArgumentNullException("userId");

            using (var command = new NpgsqlCommand("get_user_sell_triggers"))
            {
                command.CommandType = CommandType.StoredProcedure;

                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Integer,
                    Value = userId
                });

                return ExecuteGetTriggerCommand(command, GetDatabaseInstanceFromUserId(userId), userId);
            }
        }

        /// <summary>
        /// Returns a list containing all buy and sell triggers belonging to all DTS users.
        /// </summary>
        /// <returns>A list of Trigger objects or an empty list if no triggers exist in the DTS.</returns>
        public IList<Trigger> GetAllTriggers()
        {
            List<Trigger> allTriggers = new List<Trigger>();
            foreach (string connectionString in _connectionStrings)
            {
                IList<Trigger> localTriggers;
                using (var command = new NpgsqlCommand("get_all_triggers"))
                {
                    command.CommandType = CommandType.StoredProcedure;
                    localTriggers = ExecuteGetTriggerCommand(command, connectionString, null);
                }

                SetUserIdStringsOfStartupTriggers(localTriggers, connectionString);
                allTriggers.AddRange(localTriggers);
            }
            return allTriggers;
        }

        #endregion

        #region Helper Methods

        private string GetDatabaseInstanceFromUserId(string userId)
        {
            int dbIndex = userId.Select(c => (int) c).Sum() % _connectionStrings.Count;
            return _connectionStrings[dbIndex];
        }

        private void SetUserIdStringsOfStartupTriggers(IEnumerable<Trigger> triggers, string connectionString)
        {
            using (var connection = new NpgsqlConnection(connectionString))
            {
                connection.Open();
                foreach (var trigger in triggers)
                {
                    trigger.UserId = GetUserIdStringFromDbId(trigger.UserDbId, connection);
                }
            }
        }

        private string GetUserIdStringFromDbId(int userDbId, NpgsqlConnection connection)
        {
            using (var command = new NpgsqlCommand("get_user_by_id"))
            {
                command.CommandType = CommandType.StoredProcedure;
                command.Connection = connection;
                command.Parameters.Add(new NpgsqlParameter
                {
                    NpgsqlDbType = NpgsqlDbType.Integer,
                    Value = userDbId
                });

                Log.DebugFormat(CultureInfo.InvariantCulture, "Executing database query \"{0}\" using connection string \"{1}\"...", command.CommandText, connection.ConnectionString);

                /* Expected record:
                 * 
                 * (0) id
                 * (1) user_id
                 * (2) balance
                 */
                using (var reader = command.ExecuteReader())
                {
                    Log.DebugFormat("Query executed successfully.");

                    if(!reader.HasRows)
                        throw new RepositoryException(String.Format(CultureInfo.InvariantCulture,
                            "Repository does not contain a user with id {0}. Connection string: {1}", userDbId, connection.ConnectionString));

                    reader.Read();
                    return reader.GetString(1);
                }
            }
        }

        private IList<Trigger> ExecuteGetTriggerCommand(NpgsqlCommand command, string connectionString, string userId)
        {
            try
            {
                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Executing database query \"{0}\" using connection string \"{1}\"...", command.CommandText, connectionString);

                using (var connection = new NpgsqlConnection(connectionString))
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
                                var trigger = GetTriggerFromRecord(reader);
                                trigger.UserId = userId;
                                triggerList.Add(trigger);
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
                    command.CommandText, connectionString), ex);
            }
            catch (NpgsqlException ex)
            {
                throw new RepositoryException(String.Format(CultureInfo.InvariantCulture,
                    "Encountered an error while attempting to execute database command \"{0}\" using connection string \"{1}\".",
                    command.CommandText, connectionString), ex);
            }
        }

        /* Expected record:
         *      (0) id 
         *      (1) user_id
         *      (2) uid
         *      (3) stock
         *      (4) type
         *      (5) trigger_price
         *      (6) num_shares
         *      (7) created_at
         */
        private static Trigger GetTriggerFromRecord(IDataRecord reader)
        {
            return new Trigger
            {
                Id = reader.GetInt32(0),
                UserId = reader.GetString(1),
                UserDbId = reader.GetInt32(2),
                StockSymbol = reader.GetString(3),
                Type = (TriggerType) reader.GetValue(4),
                TriggerPrice = reader.GetDecimal(5),
                NumberOfShares = reader.GetInt32(6)
            };
        }

        #endregion
    }
}
