import { useSettingsStore } from '../stores/settingsStore';

export function formatCurrency(amount: number, symbol?: string): string {
    const { currencySymbol } = useSettingsStore.getState();
    const s = symbol || currencySymbol;
    return `${s} ${amount.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

export function formatDate(dateString: string, format: 'short' | 'long' | 'time' = 'short'): string {
    const date = new Date(dateString);

    switch (format) {
        case 'long':
            return date.toLocaleDateString('en-US', {
                weekday: 'long',
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            });
        case 'time':
            return date.toLocaleTimeString('en-US', {
                hour: '2-digit',
                minute: '2-digit'
            });
        case 'short':
        default:
            return date.toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric',
                year: 'numeric'
            });
    }
}

export function formatDistance(km: number): string {
    if (km < 1) {
        return `${Math.round(km * 1000)} m`;
    }
    return `${km.toFixed(1)} km`;
}

export function formatDuration(minutes: number): string {
    if (minutes < 60) {
        return `${minutes} min`;
    }
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}h ${mins}m`;
}

export function formatPhoneNumber(phone: string): string {
    // Simple formatting for Zambian numbers
    const cleaned = phone.replace(/\D/g, '');
    if (cleaned.startsWith('260')) {
        return `+${cleaned.slice(0, 3)} ${cleaned.slice(3, 5)} ${cleaned.slice(5, 8)} ${cleaned.slice(8)}`;
    }
    return phone;
}

export function truncateText(text: string, maxLength: number): string {
    if (text.length <= maxLength) return text;
    return text.slice(0, maxLength - 3) + '...';
}

export function getInitials(name: string): string {
    return name
        .split(' ')
        .map((n) => n.charAt(0).toUpperCase())
        .slice(0, 2)
        .join('');
}

export function getRelativeTime(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const diff = now.getTime() - date.getTime();

    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 1) return 'Just now';
    if (minutes < 60) return `${minutes}m ago`;
    if (hours < 24) return `${hours}h ago`;
    if (days < 7) return `${days}d ago`;

    return formatDate(dateString, 'short');
}
