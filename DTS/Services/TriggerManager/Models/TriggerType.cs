
using NpgsqlTypes;

namespace TriggerManager.Models
{
    public enum TriggerType
    {
        [EnumLabel("buy")]
        Buy,
        [EnumLabel("sell")]
        Sell
    }
}
