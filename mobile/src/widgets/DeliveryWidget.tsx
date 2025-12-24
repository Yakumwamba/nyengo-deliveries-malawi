import React from 'react';
import { FlexWidget, TextWidget } from 'react-native-android-widget';

// Theme colors matching the app - using hex strings for widget compatibility
const COLORS = {
  primary: '#FFCC00',
  secondary: '#1A1A1A',
  secondaryLight: '#2D2D2D',
  accent: '#FF3B30',
  success: '#34C759',
  successLight: '#34C75930',
  warning: '#FF9500',
  white: '#FFFFFF',
  textMuted: '#999999',
};

interface DeliveryWidgetProps {
  isOnline: boolean;
  activeDelivery: {
    orderNumber: string;
    destination: string;
    eta: string;
  } | null;
  stats: {
    todaysOrders: number;
    todaysEarnings: string;
    pending: number;
  };
}

/**
 * Small Widget (2x1) - Status Toggle + Quick Info
 */
export function SmallDeliveryWidget({ isOnline, stats }: DeliveryWidgetProps) {
  return (
    <FlexWidget
      style={{
        height: 'match_parent',
        width: 'match_parent',
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        backgroundColor: COLORS.secondary,
        borderRadius: 16,
        padding: 12,
      }}
      clickAction="OPEN_APP"
    >
      {/* Status Indicator */}
      <FlexWidget
        style={{
          flexDirection: 'row',
          alignItems: 'center',
        }}
      >
        <FlexWidget
          style={{
            width: 12,
            height: 12,
            borderRadius: 6,
            backgroundColor: isOnline ? COLORS.success : COLORS.accent,
            marginRight: 8,
          }}
        />
        <TextWidget
          text={isOnline ? 'Online' : 'Offline'}
          style={{
            fontSize: 14,
            fontWeight: '600',
            color: COLORS.white,
          }}
        />
      </FlexWidget>

      {/* Quick Stats */}
      <FlexWidget
        style={{
          alignItems: 'flex-end',
        }}
      >
        <TextWidget
          text={`${stats.todaysOrders} orders`}
          style={{
            fontSize: 12,
            color: COLORS.primary,
            fontWeight: '600',
          }}
        />
        <TextWidget
          text={stats.todaysEarnings}
          style={{
            fontSize: 14,
            fontWeight: 'bold',
            color: COLORS.white,
          }}
        />
      </FlexWidget>
    </FlexWidget>
  );
}

/**
 * Medium Widget (4x1) - Active Delivery Info
 */
export function MediumDeliveryWidget({ isOnline, activeDelivery, stats }: DeliveryWidgetProps) {
  return (
    <FlexWidget
      style={{
        height: 'match_parent',
        width: 'match_parent',
        flexDirection: 'column',
        backgroundColor: COLORS.secondary,
        borderRadius: 20,
        padding: 16,
      }}
      clickAction="OPEN_APP"
    >
      {/* Header */}
      <FlexWidget
        style={{
          flexDirection: 'row',
          justifyContent: 'space-between',
          alignItems: 'center',
          width: 'match_parent',
          marginBottom: 12,
        }}
      >
        {/* Logo & Status */}
        <FlexWidget
          style={{
            flexDirection: 'row',
            alignItems: 'center',
          }}
        >
          <FlexWidget
            style={{
              width: 32,
              height: 32,
              backgroundColor: COLORS.primary,
              borderRadius: 8,
              justifyContent: 'center',
              alignItems: 'center',
              marginRight: 8,
            }}
          >
            <TextWidget
              text="N"
              style={{
                fontSize: 18,
                fontWeight: 'bold',
                color: COLORS.secondary,
              }}
            />
          </FlexWidget>
          <FlexWidget
            style={{
              flexDirection: 'row',
              alignItems: 'center',
              backgroundColor: isOnline ? COLORS.successLight : '#FF3B3030',
              paddingHorizontal: 8,
              paddingVertical: 4,
              borderRadius: 12,
            }}
          >
            <FlexWidget
              style={{
                width: 6,
                height: 6,
                borderRadius: 3,
                backgroundColor: isOnline ? COLORS.success : COLORS.accent,
                marginRight: 4,
              }}
            />
            <TextWidget
              text={isOnline ? 'Online' : 'Offline'}
              style={{
                fontSize: 11,
                fontWeight: '600',
                color: isOnline ? COLORS.success : COLORS.accent,
              }}
            />
          </FlexWidget>
        </FlexWidget>

        {/* Stats */}
        <FlexWidget
          style={{
            flexDirection: 'row',
          }}
        >
          <FlexWidget style={{ alignItems: 'center', marginRight: 16 }}>
            <TextWidget
              text={`${stats.todaysOrders}`}
              style={{
                fontSize: 16,
                fontWeight: 'bold',
                color: COLORS.white,
              }}
            />
            <TextWidget
              text="Orders"
              style={{
                fontSize: 10,
                color: COLORS.textMuted,
              }}
            />
          </FlexWidget>
          <FlexWidget style={{ alignItems: 'center' }}>
            <TextWidget
              text={stats.todaysEarnings}
              style={{
                fontSize: 16,
                fontWeight: 'bold',
                color: COLORS.primary,
              }}
            />
            <TextWidget
              text="Earned"
              style={{
                fontSize: 10,
                color: COLORS.textMuted,
              }}
            />
          </FlexWidget>
        </FlexWidget>
      </FlexWidget>

      {/* Active Delivery or Waiting */}
      {activeDelivery ? (
        <FlexWidget
          style={{
            backgroundColor: '#FFCC0030',
            borderRadius: 12,
            padding: 12,
            width: 'match_parent',
            flexDirection: 'row',
            alignItems: 'center',
            justifyContent: 'space-between',
          }}
          clickAction="NAVIGATE"
          clickActionData={{ orderId: activeDelivery.orderNumber }}
        >
          <FlexWidget style={{ flex: 1 }}>
            <TextWidget
              text="Active Delivery"
              style={{
                fontSize: 10,
                color: COLORS.primary,
                fontWeight: '600',
              }}
            />
            <TextWidget
              text={activeDelivery.destination}
              style={{
                fontSize: 14,
                fontWeight: '600',
                color: COLORS.white,
              }}
              truncate="END"
              maxLines={1}
            />
          </FlexWidget>
          <FlexWidget
            style={{
              backgroundColor: COLORS.primary,
              paddingHorizontal: 10,
              paddingVertical: 6,
              borderRadius: 8,
            }}
          >
            <TextWidget
              text={activeDelivery.eta}
              style={{
                fontSize: 12,
                fontWeight: 'bold',
                color: COLORS.secondary,
              }}
            />
          </FlexWidget>
        </FlexWidget>
      ) : (
        <FlexWidget
          style={{
            backgroundColor: COLORS.secondaryLight,
            borderRadius: 12,
            padding: 12,
            width: 'match_parent',
            alignItems: 'center',
          }}
        >
          <TextWidget
            text={isOnline ? '🔍 Waiting for orders...' : 'Go online to receive orders'}
            style={{
              fontSize: 13,
              color: COLORS.textMuted,
            }}
          />
        </FlexWidget>
      )}
    </FlexWidget>
  );
}

/**
 * Large Widget (4x2) - Full Dashboard
 */
export function LargeDeliveryWidget({ isOnline, activeDelivery, stats }: DeliveryWidgetProps) {
  return (
    <FlexWidget
      style={{
        height: 'match_parent',
        width: 'match_parent',
        flexDirection: 'column',
        backgroundColor: COLORS.secondary,
        borderRadius: 24,
        padding: 16,
      }}
      clickAction="OPEN_APP"
    >
      {/* Header */}
      <FlexWidget
        style={{
          flexDirection: 'row',
          justifyContent: 'space-between',
          alignItems: 'center',
          width: 'match_parent',
          marginBottom: 16,
        }}
      >
        {/* Logo */}
        <FlexWidget
          style={{
            flexDirection: 'row',
            alignItems: 'center',
          }}
        >
          <FlexWidget
            style={{
              width: 40,
              height: 40,
              backgroundColor: COLORS.primary,
              borderRadius: 12,
              justifyContent: 'center',
              alignItems: 'center',
              marginRight: 10,
            }}
          >
            <TextWidget
              text="N"
              style={{
                fontSize: 22,
                fontWeight: 'bold',
                color: COLORS.secondary,
              }}
            />
          </FlexWidget>
          <FlexWidget>
            <TextWidget
              text="Nyengo"
              style={{
                fontSize: 18,
                fontWeight: 'bold',
                color: COLORS.white,
              }}
            />
            <TextWidget
              text="Deliveries"
              style={{
                fontSize: 11,
                color: COLORS.textMuted,
              }}
            />
          </FlexWidget>
        </FlexWidget>

        {/* Status Toggle */}
        <FlexWidget
          style={{
            flexDirection: 'row',
            alignItems: 'center',
            backgroundColor: isOnline ? COLORS.success : '#404040',
            paddingHorizontal: 12,
            paddingVertical: 8,
            borderRadius: 20,
          }}
          clickAction="TOGGLE_STATUS"
        >
          <FlexWidget
            style={{
              width: 8,
              height: 8,
              borderRadius: 4,
              backgroundColor: isOnline ? COLORS.white : COLORS.textMuted,
              marginRight: 6,
            }}
          />
          <TextWidget
            text={isOnline ? 'Online' : 'Offline'}
            style={{
              fontSize: 12,
              fontWeight: '600',
              color: COLORS.white,
            }}
          />
        </FlexWidget>
      </FlexWidget>

      {/* Stats Grid */}
      <FlexWidget
        style={{
          flexDirection: 'row',
          width: 'match_parent',
          marginBottom: 16,
        }}
      >
        <FlexWidget
          style={{
            flex: 1,
            backgroundColor: COLORS.secondaryLight,
            borderRadius: 12,
            padding: 10,
            alignItems: 'center',
            marginRight: 5,
          }}
        >
          <TextWidget text="📦" style={{ fontSize: 16, marginBottom: 4 }} />
          <TextWidget
            text={`${stats.todaysOrders}`}
            style={{
              fontSize: 16,
              fontWeight: 'bold',
              color: COLORS.primary,
            }}
          />
          <TextWidget
            text="Orders"
            style={{
              fontSize: 9,
              color: COLORS.textMuted,
            }}
          />
        </FlexWidget>
        <FlexWidget
          style={{
            flex: 1,
            backgroundColor: COLORS.secondaryLight,
            borderRadius: 12,
            padding: 10,
            alignItems: 'center',
            marginHorizontal: 5,
          }}
        >
          <TextWidget text="💰" style={{ fontSize: 16, marginBottom: 4 }} />
          <TextWidget
            text={stats.todaysEarnings}
            style={{
              fontSize: 16,
              fontWeight: 'bold',
              color: COLORS.success,
            }}
          />
          <TextWidget
            text="Earned"
            style={{
              fontSize: 9,
              color: COLORS.textMuted,
            }}
          />
        </FlexWidget>
        <FlexWidget
          style={{
            flex: 1,
            backgroundColor: COLORS.secondaryLight,
            borderRadius: 12,
            padding: 10,
            alignItems: 'center',
            marginLeft: 5,
          }}
        >
          <TextWidget text="⏳" style={{ fontSize: 16, marginBottom: 4 }} />
          <TextWidget
            text={`${stats.pending}`}
            style={{
              fontSize: 16,
              fontWeight: 'bold',
              color: COLORS.warning,
            }}
          />
          <TextWidget
            text="Pending"
            style={{
              fontSize: 9,
              color: COLORS.textMuted,
            }}
          />
        </FlexWidget>
      </FlexWidget>

      {/* Active Delivery Card */}
      {activeDelivery ? (
        <FlexWidget
          style={{
            backgroundColor: COLORS.primary,
            borderRadius: 16,
            padding: 14,
            width: 'match_parent',
            flexDirection: 'row',
            alignItems: 'center',
            justifyContent: 'space-between',
          }}
          clickAction="NAVIGATE"
          clickActionData={{ orderId: activeDelivery.orderNumber }}
        >
          <FlexWidget style={{ flex: 1, marginRight: 12 }}>
            <FlexWidget
              style={{
                flexDirection: 'row',
                alignItems: 'center',
                marginBottom: 4,
              }}
            >
              <TextWidget
                text="🚴 ACTIVE DELIVERY"
                style={{
                  fontSize: 10,
                  color: COLORS.secondary,
                  fontWeight: '700',
                }}
              />
            </FlexWidget>
            <TextWidget
              text={activeDelivery.orderNumber}
              style={{
                fontSize: 12,
                color: '#1A1A1ACC',
              }}
            />
            <TextWidget
              text={activeDelivery.destination}
              style={{
                fontSize: 15,
                fontWeight: 'bold',
                color: COLORS.secondary,
              }}
              truncate="END"
              maxLines={1}
            />
          </FlexWidget>
          
          <FlexWidget
            style={{
              backgroundColor: COLORS.secondary,
              paddingHorizontal: 14,
              paddingVertical: 10,
              borderRadius: 12,
              alignItems: 'center',
            }}
          >
            <TextWidget
              text={activeDelivery.eta}
              style={{
                fontSize: 18,
                fontWeight: 'bold',
                color: COLORS.primary,
              }}
            />
            <TextWidget
              text="ETA"
              style={{
                fontSize: 9,
                color: COLORS.textMuted,
                fontWeight: '600',
              }}
            />
          </FlexWidget>
        </FlexWidget>
      ) : (
        <FlexWidget
          style={{
            backgroundColor: COLORS.secondaryLight,
            borderRadius: 16,
            padding: 16,
            width: 'match_parent',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <TextWidget
            text={isOnline ? '🔍' : '🔴'}
            style={{ fontSize: 24, marginBottom: 4 }}
          />
          <TextWidget
            text={isOnline ? 'Waiting for orders...' : 'You are offline'}
            style={{
              fontSize: 14,
              color: COLORS.textMuted,
              textAlign: 'center',
            }}
          />
          {!isOnline && (
            <TextWidget
              text="Tap to go online"
              style={{
                fontSize: 12,
                color: COLORS.primary,
                marginTop: 4,
              }}
            />
          )}
        </FlexWidget>
      )}
    </FlexWidget>
  );
}

export default LargeDeliveryWidget;
