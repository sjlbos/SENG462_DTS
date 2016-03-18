using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using log4net;

namespace TriggerManager.Models
{
    /// <summary>
    /// A thread-safe class in charge of managing and firing the active buy and sell triggers in the DTS.
    /// </summary>
    public class TriggerController
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (TriggerController));

        private IDictionary<string, IList<Trigger>> _userToBuyTriggersMap;
        private IDictionary<string, IList<Trigger>> _userToSellTriggersMap;
        private readonly TriggerLogicManager _triggerLogicManager;
        private readonly ITriggerAuthority _triggerAuthority;
        private readonly object _syncLock;

        public TriggerController(ITriggerAuthority triggerAuthority)
        {
            if (triggerAuthority == null)
                throw new ArgumentNullException("triggerAuthority");

            _userToBuyTriggersMap = new Dictionary<string, IList<Trigger>>();
            _userToSellTriggersMap = new Dictionary<string, IList<Trigger>>();
            _triggerLogicManager = new TriggerLogicManager();
            _triggerAuthority = triggerAuthority;
            _syncLock = new object();
        }

        #region Public Interface

        /// <summary>
        /// Removes all triggers managed by the TriggerController and replaces them with the provided triggers. 
        /// Can be used to initialize the TriggerController's list of triggers.
        /// </summary>
        /// <param name="triggers">The list of trigger objects the TriggerController should oversee.</param>
        public void UpdateTriggers(IList<Trigger> triggers)
        {
            if (triggers == null)
                throw new ArgumentNullException("triggers");

            lock (_syncLock)
            {
                // Drop all existing triggers
                _triggerLogicManager.Reset();
                _userToBuyTriggersMap = new Dictionary<string, IList<Trigger>>();
                _userToSellTriggersMap = new Dictionary<string, IList<Trigger>>();

                // Add new triggers
                foreach (var trigger in triggers)
                {
                    _triggerLogicManager.AddTrigger(trigger);
                    if (trigger.Type == TriggerType.Buy)
                    {
                        AddBuyTrigger(trigger);
                    }
                    else
                    {
                        AddSellTrigger(trigger);
                    }
                }
            }
        }

        /// <summary>
        /// Removes all buy triggers belonging to the spcified user from the TriggerController and replaces them
        /// with the provided list of triggers.
        /// </summary>
        /// <param name="userId">The user ID of the buy triggers' owner.</param>
        /// <param name="triggers">The user's new buy triggers.</param>
        public void UpdateBuyTriggersForUser(string userId, IList<Trigger> triggers)
        {
            lock (_syncLock)
            {
                if (_userToBuyTriggersMap.ContainsKey(userId))
                {
                    UpdateBuyTriggersForExistingUser(userId, triggers);
                }
                else
                {
                    AddBuyTriggersForNewUser(userId, triggers);
                }
            }
        }

        /// <summary>
        /// Removes all sell triggers belonging to the spcified user from the TriggerController and replaces them
        /// with the provided list of triggers.
        /// </summary>
        /// <param name="userId">The user ID of the sell triggers' owner.</param>
        /// <param name="triggers">The users's new sell triggers.</param>
        public void UpdateSellTriggersForUser(string userId, IList<Trigger> triggers)
        {
            lock (_syncLock)
            {
                if (_userToSellTriggersMap.ContainsKey(userId))
                {
                    UpdateSellTriggersForExistingUser(userId, triggers);
                }
                else
                {
                    AddSellTriggersForNewUser(userId, triggers);
                }
            }
        }

        /// <summary>
        /// Handles the firing of buy and sell triggers when the provided stock reaches the provided price.
        /// </summary>
        /// <param name="stockSymbol">A 3-character stock symbol.</param>
        /// <param name="price">The current price of the stock represented by the stock symbol.</param>
        public void HandleStockPriceUpdate(string stockSymbol, decimal price)
        {
            if (stockSymbol == null)
                throw new ArgumentNullException("stockSymbol");

            IList<Trigger> sellTriggersToFire;
            IList<Trigger> buyTriggersToFire;

            lock (_syncLock)
            {
                sellTriggersToFire = _triggerLogicManager.GetFiredSellTriggersAtStockPrice(stockSymbol, price);
                buyTriggersToFire = _triggerLogicManager.GetFiredBuyTriggersAtStockPrice(stockSymbol, price);
            }

            FireReadyBuyTriggers(buyTriggersToFire);
            FireReadySellTriggers(sellTriggersToFire);
        }

        /// <summary>
        /// The complete list of stock symbols the TriggerController is currently monitoring.
        /// </summary>
        public IEnumerable<string> StockList
        {
            get { return _triggerLogicManager.StockList; }
        }  

        #endregion

        #region Helper Methods

        private void AddBuyTrigger(Trigger trigger)
        {
            if (!_userToBuyTriggersMap.ContainsKey(trigger.UserId))
            {
                _userToBuyTriggersMap.Add(trigger.UserId, new List<Trigger>());
            }
            _userToBuyTriggersMap[trigger.UserId].Add(trigger);
        }

        private void AddSellTrigger(Trigger trigger)
        {
            if (!_userToSellTriggersMap.ContainsKey(trigger.UserId))
            {
                _userToSellTriggersMap.Add(trigger.UserId, new List<Trigger>());
            }
            _userToSellTriggersMap[trigger.UserId].Add(trigger);
        }

        private void UpdateBuyTriggersForExistingUser(string userId, IList<Trigger> triggers)
        {
            var currentBuyTriggers = _userToBuyTriggersMap[userId];
            foreach (var trigger in currentBuyTriggers)
            {
                _triggerLogicManager.RemoveTrigger(trigger);
            }
            if (triggers == null)
            {
                _userToBuyTriggersMap.Remove(userId);
            }
            else
            {
                _userToBuyTriggersMap[userId] = triggers;
            }
        }

        private void AddBuyTriggersForNewUser(string userId, IList<Trigger> triggers)
        {
            if (triggers == null)
                return;
            _userToBuyTriggersMap.Add(userId, triggers);
            foreach (var trigger in triggers)
            {
                _triggerLogicManager.AddTrigger(trigger);
            }
        }

        private void UpdateSellTriggersForExistingUser(string userId, IList<Trigger> triggers)
        {
            var currentSellTriggers = _userToSellTriggersMap[userId];
            foreach (var trigger in currentSellTriggers)
            {
                _triggerLogicManager.RemoveTrigger(trigger);
            }
            if (triggers == null)
            {
                _userToSellTriggersMap.Remove(userId);
            }
            else
            {
                _userToSellTriggersMap[userId] = triggers;
            }
        }

        private void AddSellTriggersForNewUser(string userId, IList<Trigger> triggers)
        {
            if (triggers == null)
                return;
            _userToSellTriggersMap.Add(userId, triggers);
            foreach (var trigger in triggers)
            {
                _triggerLogicManager.AddTrigger(trigger);
            }
        }

        private void FireReadySellTriggers(IList<Trigger> sellTriggers)
        {
            if (sellTriggers == null || sellTriggers.Count == 0)
                return;
            foreach (var trigger in sellTriggers)
            {
                Trigger localTrigger = trigger; // copy to local variable to avoid using iterator variable in lambda 
                Task.Run(() => ExecuteTrigger(localTrigger));
            }
        }

        private void FireReadyBuyTriggers(IList<Trigger> buyTriggers)
        {
            if (buyTriggers == null || buyTriggers.Count == 0)
                return;

            foreach (var trigger in buyTriggers)
            {
                Trigger localTrigger = trigger; // copy to local variable to avoid using iterator variable in lambda 
                Task.Run(() => ExecuteTrigger(localTrigger));
            }
        }

        private void ExecuteTrigger(Trigger trigger)
        {
            try
            {
                _triggerAuthority.ExecuteTrigger(trigger);
            }
            catch (TriggerCommitException ex)
            {
                Log.Error(ex);
            }
        }
        #endregion
    }
}
