
using TriggerManager.Models;

namespace TriggerManager
{
    public interface ITriggerAuthority
    {
        void ExecuteTrigger(Trigger trigger);
    }
}
