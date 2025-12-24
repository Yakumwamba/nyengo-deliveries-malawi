export function formatCurrencySimple(amount: number, symbol: string = 'K'): string {
    return `${symbol} ${amount.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

export function parseCurrency(value: string): number {
    const cleaned = value.replace(/[^0-9.-]/g, '');
    return parseFloat(cleaned) || 0;
}

export const CURRENCY_CONFIG: Record<string, { symbol: string; locale: string }> = {
    ZMW: { symbol: 'K', locale: 'en-ZM' },
    USD: { symbol: '$', locale: 'en-US' },
    ZAR: { symbol: 'R', locale: 'en-ZA' },
    KES: { symbol: 'KSh', locale: 'en-KE' },
    NGN: { symbol: '₦', locale: 'en-NG' },
    GHS: { symbol: 'GH₵', locale: 'en-GH' },
};

export function getCurrencySymbol(code: string): string {
    return CURRENCY_CONFIG[code]?.symbol || code;
}
