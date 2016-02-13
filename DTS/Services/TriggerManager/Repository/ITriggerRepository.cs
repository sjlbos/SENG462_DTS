
using System.Collections.Generic;
using TriggerManager.Models;

namespace TriggerManager.Repository
{
    public interface ITriggerRepository
    {
        IList<Trigger> GetBuyTriggersForUser(int userId);
        IList<Trigger> GetSellTriggersForUser(int userId);
        IList<Trigger> GetAllTriggers();
    }
}
