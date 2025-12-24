/**
 * Nyengo Deliveries - Yandex/Yango Inspired Theme
 * Professional Design System with Yellow & Red/Black Theme
 */

// ============================================
// COLORS - Yandex/Yango Inspired Palette
// ============================================
export const COLORS = {
    // Primary - Signature Yandex Yellow
    primary: '#FFCC00',
    primaryDark: '#E6B800',
    primaryLight: '#FFE066',
    primaryMuted: '#FFF5CC',

    // Secondary - Deep Black/Charcoal (Yango style)
    secondary: '#1A1A1A',
    secondaryLight: '#2D2D2D',
    secondaryMuted: '#404040',

    // Accent - Vibrant Red (Yango accent)
    accent: '#FF3B30',
    accentDark: '#E6352B',
    accentLight: '#FF6B63',
    accentMuted: '#FFEBEA',

    // Semantic Colors
    success: '#34C759',
    successLight: '#E8F9ED',
    warning: '#FF9500',
    warningLight: '#FFF3E0',
    danger: '#FF3B30',
    dangerLight: '#FFEBEA',
    info: '#007AFF',
    infoLight: '#E5F1FF',

    // Neutrals
    white: '#FFFFFF',
    black: '#000000',
    background: '#F8F8F8',
    backgroundDark: '#1A1A1A',
    surface: '#FFFFFF',
    surfaceElevated: '#FFFFFF',

    // Text Colors
    text: '#1A1A1A',
    textSecondary: '#666666',
    textMuted: '#999999',
    textLight: '#CCCCCC',
    textInverse: '#FFFFFF',

    // Border & Dividers
    border: '#E5E5E5',
    borderLight: '#F0F0F0',
    divider: '#EEEEEE',

    // Overlay
    overlay: 'rgba(0, 0, 0, 0.5)',
    overlayLight: 'rgba(0, 0, 0, 0.3)',

    // Gradient combinations
    gradientYellow: ['#FFCC00', '#FFE066'],
    gradientRed: ['#FF3B30', '#FF6B63'],
    gradientDark: ['#1A1A1A', '#2D2D2D'],
};

// ============================================
// TYPOGRAPHY - Modern, Clean Font System
// ============================================
export const TYPOGRAPHY = {
    // Font Families
    fontFamily: {
        regular: 'System',
        medium: 'System',
        semibold: 'System',
        bold: 'System',
    },

    // Font Sizes
    fontSize: {
        xs: 10,
        sm: 12,
        base: 14,
        md: 16,
        lg: 18,
        xl: 20,
        '2xl': 24,
        '3xl': 28,
        '4xl': 32,
        '5xl': 40,
        '6xl': 48,
    },

    // Font Weights
    fontWeight: {
        regular: '400' as const,
        medium: '500' as const,
        semibold: '600' as const,
        bold: '700' as const,
        extrabold: '800' as const,
    },

    // Line Heights
    lineHeight: {
        tight: 1.2,
        snug: 1.375,
        normal: 1.5,
        relaxed: 1.625,
        loose: 2,
    },

    // Letter Spacing
    letterSpacing: {
        tighter: -0.5,
        tight: -0.25,
        normal: 0,
        wide: 0.25,
        wider: 0.5,
        widest: 1,
    },
};

// ============================================
// TEXT STYLES - Pre-defined Text Variations
// ============================================
export const TEXT_STYLES = {
    // Headings
    h1: {
        fontSize: TYPOGRAPHY.fontSize['4xl'],
        fontWeight: TYPOGRAPHY.fontWeight.bold,
        lineHeight: TYPOGRAPHY.fontSize['4xl'] * TYPOGRAPHY.lineHeight.tight,
        color: COLORS.text,
        letterSpacing: TYPOGRAPHY.letterSpacing.tight,
    },
    h2: {
        fontSize: TYPOGRAPHY.fontSize['3xl'],
        fontWeight: TYPOGRAPHY.fontWeight.bold,
        lineHeight: TYPOGRAPHY.fontSize['3xl'] * TYPOGRAPHY.lineHeight.tight,
        color: COLORS.text,
        letterSpacing: TYPOGRAPHY.letterSpacing.tight,
    },
    h3: {
        fontSize: TYPOGRAPHY.fontSize['2xl'],
        fontWeight: TYPOGRAPHY.fontWeight.semibold,
        lineHeight: TYPOGRAPHY.fontSize['2xl'] * TYPOGRAPHY.lineHeight.snug,
        color: COLORS.text,
    },
    h4: {
        fontSize: TYPOGRAPHY.fontSize.xl,
        fontWeight: TYPOGRAPHY.fontWeight.semibold,
        lineHeight: TYPOGRAPHY.fontSize.xl * TYPOGRAPHY.lineHeight.snug,
        color: COLORS.text,
    },
    h5: {
        fontSize: TYPOGRAPHY.fontSize.lg,
        fontWeight: TYPOGRAPHY.fontWeight.semibold,
        lineHeight: TYPOGRAPHY.fontSize.lg * TYPOGRAPHY.lineHeight.snug,
        color: COLORS.text,
    },

    // Body Text
    bodyLarge: {
        fontSize: TYPOGRAPHY.fontSize.md,
        fontWeight: TYPOGRAPHY.fontWeight.regular,
        lineHeight: TYPOGRAPHY.fontSize.md * TYPOGRAPHY.lineHeight.normal,
        color: COLORS.text,
    },
    body: {
        fontSize: TYPOGRAPHY.fontSize.base,
        fontWeight: TYPOGRAPHY.fontWeight.regular,
        lineHeight: TYPOGRAPHY.fontSize.base * TYPOGRAPHY.lineHeight.normal,
        color: COLORS.text,
    },
    bodySmall: {
        fontSize: TYPOGRAPHY.fontSize.sm,
        fontWeight: TYPOGRAPHY.fontWeight.regular,
        lineHeight: TYPOGRAPHY.fontSize.sm * TYPOGRAPHY.lineHeight.normal,
        color: COLORS.textSecondary,
    },

    // Utility Text
    caption: {
        fontSize: TYPOGRAPHY.fontSize.xs,
        fontWeight: TYPOGRAPHY.fontWeight.regular,
        lineHeight: TYPOGRAPHY.fontSize.xs * TYPOGRAPHY.lineHeight.normal,
        color: COLORS.textMuted,
        letterSpacing: TYPOGRAPHY.letterSpacing.wide,
    },
    label: {
        fontSize: TYPOGRAPHY.fontSize.sm,
        fontWeight: TYPOGRAPHY.fontWeight.semibold,
        lineHeight: TYPOGRAPHY.fontSize.sm * TYPOGRAPHY.lineHeight.normal,
        color: COLORS.textSecondary,
        textTransform: 'uppercase' as const,
        letterSpacing: TYPOGRAPHY.letterSpacing.wider,
    },
    button: {
        fontSize: TYPOGRAPHY.fontSize.md,
        fontWeight: TYPOGRAPHY.fontWeight.semibold,
        lineHeight: TYPOGRAPHY.fontSize.md * TYPOGRAPHY.lineHeight.tight,
        letterSpacing: TYPOGRAPHY.letterSpacing.wide,
    },
    buttonSmall: {
        fontSize: TYPOGRAPHY.fontSize.base,
        fontWeight: TYPOGRAPHY.fontWeight.semibold,
        lineHeight: TYPOGRAPHY.fontSize.base * TYPOGRAPHY.lineHeight.tight,
        letterSpacing: TYPOGRAPHY.letterSpacing.wide,
    },
};

// ============================================
// SPACING - Consistent Spacing Scale
// ============================================
export const SPACING = {
    xs: 4,
    sm: 8,
    md: 12,
    base: 16,
    lg: 20,
    xl: 24,
    '2xl': 32,
    '3xl': 40,
    '4xl': 48,
    '5xl': 64,
};

// ============================================
// BORDER RADIUS
// ============================================
export const RADIUS = {
    none: 0,
    sm: 4,
    md: 8,
    base: 12,
    lg: 16,
    xl: 20,
    '2xl': 24,
    full: 9999,
};

// ============================================
// SHADOWS
// ============================================
export const SHADOWS = {
    sm: {
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 1 },
        shadowOpacity: 0.05,
        shadowRadius: 2,
        elevation: 1,
    },
    md: {
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.08,
        shadowRadius: 4,
        elevation: 2,
    },
    lg: {
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.1,
        shadowRadius: 8,
        elevation: 4,
    },
    xl: {
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 8 },
        shadowOpacity: 0.15,
        shadowRadius: 16,
        elevation: 8,
    },
    // Colored shadows for CTAs
    primary: {
        shadowColor: COLORS.primary,
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.4,
        shadowRadius: 8,
        elevation: 4,
    },
    accent: {
        shadowColor: COLORS.accent,
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.3,
        shadowRadius: 8,
        elevation: 4,
    },
};

// ============================================
// STATUS COLORS
// ============================================
export const STATUS_COLORS = {
    pending: {
        bg: COLORS.warningLight,
        text: COLORS.warning,
        icon: COLORS.warning,
    },
    accepted: {
        bg: COLORS.primaryMuted,
        text: COLORS.primaryDark,
        icon: COLORS.primary,
    },
    in_transit: {
        bg: COLORS.infoLight,
        text: COLORS.info,
        icon: COLORS.info,
    },
    picked_up: {
        bg: COLORS.infoLight,
        text: COLORS.info,
        icon: COLORS.info,
    },
    delivered: {
        bg: COLORS.successLight,
        text: COLORS.success,
        icon: COLORS.success,
    },
    cancelled: {
        bg: COLORS.dangerLight,
        text: COLORS.danger,
        icon: COLORS.danger,
    },
    failed: {
        bg: COLORS.dangerLight,
        text: COLORS.danger,
        icon: COLORS.danger,
    },
};

// ============================================
// COMPONENT STYLES
// ============================================
export const COMPONENT_STYLES = {
    // Cards
    card: {
        backgroundColor: COLORS.surface,
        borderRadius: RADIUS.lg,
        padding: SPACING.base,
        ...SHADOWS.md,
    },
    cardElevated: {
        backgroundColor: COLORS.surfaceElevated,
        borderRadius: RADIUS.xl,
        padding: SPACING.lg,
        ...SHADOWS.lg,
    },

    // Buttons
    buttonPrimary: {
        backgroundColor: COLORS.primary,
        borderRadius: RADIUS.base,
        height: 56,
        justifyContent: 'center' as const,
        alignItems: 'center' as const,
        ...SHADOWS.primary,
    },
    buttonSecondary: {
        backgroundColor: COLORS.secondary,
        borderRadius: RADIUS.base,
        height: 56,
        justifyContent: 'center' as const,
        alignItems: 'center' as const,
        ...SHADOWS.lg,
    },
    buttonAccent: {
        backgroundColor: COLORS.accent,
        borderRadius: RADIUS.base,
        height: 56,
        justifyContent: 'center' as const,
        alignItems: 'center' as const,
        ...SHADOWS.accent,
    },
    buttonOutline: {
        backgroundColor: 'transparent',
        borderRadius: RADIUS.base,
        height: 56,
        justifyContent: 'center' as const,
        alignItems: 'center' as const,
        borderWidth: 2,
        borderColor: COLORS.primary,
    },

    // Inputs
    input: {
        backgroundColor: COLORS.background,
        borderRadius: RADIUS.base,
        height: 56,
        paddingHorizontal: SPACING.base,
        fontSize: TYPOGRAPHY.fontSize.md,
        color: COLORS.text,
    },
    inputFocused: {
        backgroundColor: COLORS.white,
        borderWidth: 2,
        borderColor: COLORS.primary,
    },

    // Header
    header: {
        backgroundColor: COLORS.secondary,
        paddingTop: 60,
        paddingBottom: SPACING.lg,
        paddingHorizontal: SPACING.lg,
    },
    headerLight: {
        backgroundColor: COLORS.primary,
        paddingTop: 60,
        paddingBottom: SPACING.lg,
        paddingHorizontal: SPACING.lg,
    },
};

// ============================================
// ICON SIZES
// ============================================
export const ICON_SIZES = {
    xs: 16,
    sm: 20,
    md: 24,
    lg: 28,
    xl: 32,
    '2xl': 40,
    '3xl': 48,
    '4xl': 64,
};

// Default export for convenience
const theme = {
    COLORS,
    TYPOGRAPHY,
    TEXT_STYLES,
    SPACING,
    RADIUS,
    SHADOWS,
    STATUS_COLORS,
    COMPONENT_STYLES,
    ICON_SIZES,
};

export default theme;
