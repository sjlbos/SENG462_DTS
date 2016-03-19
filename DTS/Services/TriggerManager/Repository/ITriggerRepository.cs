
using System.Collections.Generic;
using TriggerManager.Models;

namespace TriggerManager.Repository
{
    public interface ITriggerRepository
    {
        IList<Trigger> GetBuyTriggersForUser(int userDbId, string userId);
        IList<Trigger> GetSellTriggersForUser(int userDbId, string userId);
        IList<Trigger> GetAllTriggers();
    }
}
