import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  Alert,
  ActivityIndicator,
} from 'react-native';
import { Link, router } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { authService } from '../../src/services/auth';
import { COLORS, TYPOGRAPHY, SPACING, RADIUS, SHADOWS } from '../../src/constants/theme';

export default function RegisterScreen() {
  const [form, setForm] = useState({
    email: '',
    password: '',
    confirmPassword: '',
    companyName: '',
    ownerName: '',
    phone: '',
    address: '',
    city: '',
  });
  const [isLoading, setIsLoading] = useState(false);
  const [step, setStep] = useState(1);
  const [focusedInput, setFocusedInput] = useState<string | null>(null);

  const updateForm = (key: string, value: string) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  };

  const handleRegister = async () => {
    if (!form.email || !form.password || !form.companyName || !form.ownerName || !form.phone) {
      Alert.alert('Error', 'Please fill in all required fields');
      return;
    }
    if (form.password !== form.confirmPassword) {
      Alert.alert('Error', 'Passwords do not match');
      return;
    }

    setIsLoading(true);
    try {
      await authService.register({
        email: form.email,
        password: form.password,
        companyName: form.companyName,
        ownerName: form.ownerName,
        phone: form.phone,
        address: form.address,
        city: form.city,
        country: 'Zambia',
        serviceAreas: [form.city],
        vehicleTypes: ['motorcycle', 'car'],
      });
      Alert.alert('Success', 'Registration successful! Please login.', [
        { text: 'OK', onPress: () => router.replace('/(auth)/login') },
      ]);
    } catch (error: any) {
      Alert.alert('Registration Failed', error.message);
    } finally {
      setIsLoading(false);
    }
  };

  const canProceedStep1 = form.email && form.password && form.confirmPassword;
  const canProceedStep2 = form.companyName && form.ownerName && form.phone && form.city;

  return (
    <ScrollView 
      style={styles.container} 
      contentContainerStyle={styles.content}
      showsVerticalScrollIndicator={false}
    >
      {/* Header */}
      <View style={styles.header}>
        <Link href="/(auth)/login" asChild>
          <TouchableOpacity style={styles.backButton}>
            <Ionicons name="arrow-back" size={24} color={COLORS.text} />
          </TouchableOpacity>
        </Link>
        
        <View style={styles.headerContent}>
          <Text style={styles.title}>Create Account</Text>
          <Text style={styles.subtitle}>Join Nyengo Deliveries</Text>
        </View>

        {/* Progress Indicator */}
        <View style={styles.progressContainer}>
          <View style={styles.progressBar}>
            <View style={[styles.progressFill, { width: step === 1 ? '50%' : '100%' }]} />
          </View>
          <Text style={styles.progressText}>Step {step} of 2</Text>
        </View>
      </View>

      {/* Form */}
      <View style={styles.form}>
        {step === 1 ? (
          <>
            <Text style={styles.sectionTitle}>
              <Ionicons name="person-circle" size={18} color={COLORS.primary} /> Account Details
            </Text>
            
            <FormInput
              icon="mail-outline"
              placeholder="Email *"
              value={form.email}
              onChangeText={(v) => updateForm('email', v)}
              keyboardType="email-address"
              autoCapitalize="none"
              focused={focusedInput === 'email'}
              onFocus={() => setFocusedInput('email')}
              onBlur={() => setFocusedInput(null)}
            />
            <FormInput
              icon="lock-closed-outline"
              placeholder="Password *"
              value={form.password}
              onChangeText={(v) => updateForm('password', v)}
              secureTextEntry
              focused={focusedInput === 'password'}
              onFocus={() => setFocusedInput('password')}
              onBlur={() => setFocusedInput(null)}
            />
            <FormInput
              icon="shield-checkmark-outline"
              placeholder="Confirm Password *"
              value={form.confirmPassword}
              onChangeText={(v) => updateForm('confirmPassword', v)}
              secureTextEntry
              focused={focusedInput === 'confirmPassword'}
              onFocus={() => setFocusedInput('confirmPassword')}
              onBlur={() => setFocusedInput(null)}
            />

            <TouchableOpacity 
              style={[styles.nextButton, !canProceedStep1 && styles.buttonDisabled]}
              onPress={() => setStep(2)}
              disabled={!canProceedStep1}
            >
              <Text style={styles.nextButtonText}>Continue</Text>
              <Ionicons name="arrow-forward" size={20} color={COLORS.secondary} />
            </TouchableOpacity>
          </>
        ) : (
          <>
            <Text style={styles.sectionTitle}>
              <Ionicons name="business" size={18} color={COLORS.primary} /> Business Details
            </Text>
            
            <FormInput
              icon="storefront-outline"
              placeholder="Company Name *"
              value={form.companyName}
              onChangeText={(v) => updateForm('companyName', v)}
              focused={focusedInput === 'companyName'}
              onFocus={() => setFocusedInput('companyName')}
              onBlur={() => setFocusedInput(null)}
            />
            <FormInput
              icon="person-outline"
              placeholder="Owner Name *"
              value={form.ownerName}
              onChangeText={(v) => updateForm('ownerName', v)}
              focused={focusedInput === 'ownerName'}
              onFocus={() => setFocusedInput('ownerName')}
              onBlur={() => setFocusedInput(null)}
            />
            <FormInput
              icon="call-outline"
              placeholder="Phone Number *"
              value={form.phone}
              onChangeText={(v) => updateForm('phone', v)}
              keyboardType="phone-pad"
              focused={focusedInput === 'phone'}
              onFocus={() => setFocusedInput('phone')}
              onBlur={() => setFocusedInput(null)}
            />
            <FormInput
              icon="location-outline"
              placeholder="Business Address"
              value={form.address}
              onChangeText={(v) => updateForm('address', v)}
              focused={focusedInput === 'address'}
              onFocus={() => setFocusedInput('address')}
              onBlur={() => setFocusedInput(null)}
            />
            <FormInput
              icon="map-outline"
              placeholder="City *"
              value={form.city}
              onChangeText={(v) => updateForm('city', v)}
              focused={focusedInput === 'city'}
              onFocus={() => setFocusedInput('city')}
              onBlur={() => setFocusedInput(null)}
            />

            <View style={styles.buttonRow}>
              <TouchableOpacity 
                style={styles.backStepButton}
                onPress={() => setStep(1)}
              >
                <Ionicons name="arrow-back" size={20} color={COLORS.text} />
                <Text style={styles.backStepButtonText}>Back</Text>
              </TouchableOpacity>
              
              <TouchableOpacity 
                style={[styles.registerButton, (!canProceedStep2 || isLoading) && styles.buttonDisabled]}
                onPress={handleRegister} 
                disabled={!canProceedStep2 || isLoading}
              >
                {isLoading ? (
                  <ActivityIndicator color={COLORS.secondary} />
                ) : (
                  <>
                    <Text style={styles.registerButtonText}>Register</Text>
                    <Ionicons name="checkmark-circle" size={20} color={COLORS.secondary} />
                  </>
                )}
              </TouchableOpacity>
            </View>
          </>
        )}
      </View>

      {/* Login Link */}
      <View style={styles.loginContainer}>
        <Text style={styles.loginText}>Already have an account? </Text>
        <Link href="/(auth)/login" asChild>
          <TouchableOpacity>
            <Text style={styles.loginLink}>Login</Text>
          </TouchableOpacity>
        </Link>
      </View>
    </ScrollView>
  );
}

function FormInput({
  icon,
  placeholder,
  value,
  onChangeText,
  keyboardType,
  autoCapitalize,
  secureTextEntry,
  focused,
  onFocus,
  onBlur,
}: {
  icon: string;
  placeholder: string;
  value: string;
  onChangeText: (text: string) => void;
  keyboardType?: 'default' | 'email-address' | 'phone-pad';
  autoCapitalize?: 'none' | 'sentences' | 'words' | 'characters';
  secureTextEntry?: boolean;
  focused: boolean;
  onFocus: () => void;
  onBlur: () => void;
}) {
  return (
    <View style={[styles.inputContainer, focused && styles.inputContainerFocused]}>
      <View style={styles.inputIcon}>
        <Ionicons 
          name={icon as any} 
          size={20} 
          color={focused ? COLORS.primary : COLORS.textMuted} 
        />
      </View>
      <TextInput
        style={styles.input}
        placeholder={placeholder}
        value={value}
        onChangeText={onChangeText}
        keyboardType={keyboardType}
        autoCapitalize={autoCapitalize}
        secureTextEntry={secureTextEntry}
        placeholderTextColor={COLORS.textMuted}
        onFocus={onFocus}
        onBlur={onBlur}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { 
    flex: 1, 
    backgroundColor: COLORS.white,
  },
  content: { 
    padding: SPACING.xl,
    paddingTop: 0,
  },
  
  // Header
  header: { 
    paddingTop: 60,
    marginBottom: SPACING.xl,
  },
  backButton: { 
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: COLORS.background,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: SPACING.lg,
  },
  headerContent: {
    marginBottom: SPACING.lg,
  },
  title: { 
    fontSize: TYPOGRAPHY.fontSize['3xl'], 
    fontWeight: TYPOGRAPHY.fontWeight.bold, 
    color: COLORS.text,
    letterSpacing: -0.5,
  },
  subtitle: { 
    fontSize: TYPOGRAPHY.fontSize.md, 
    color: COLORS.textSecondary, 
    marginTop: 4,
  },
  progressContainer: {
    gap: SPACING.sm,
  },
  progressBar: {
    height: 6,
    backgroundColor: COLORS.border,
    borderRadius: 3,
    overflow: 'hidden',
  },
  progressFill: {
    height: '100%',
    backgroundColor: COLORS.primary,
    borderRadius: 3,
  },
  progressText: {
    fontSize: TYPOGRAPHY.fontSize.sm,
    color: COLORS.textMuted,
    fontWeight: TYPOGRAPHY.fontWeight.medium,
  },
  
  // Form
  form: { 
    marginBottom: SPACING.xl,
  },
  sectionTitle: { 
    fontSize: TYPOGRAPHY.fontSize.md, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold, 
    color: COLORS.text,
    marginBottom: SPACING.md,
    marginTop: SPACING.sm,
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: COLORS.background,
    borderRadius: RADIUS.lg,
    marginBottom: SPACING.md,
    height: 56,
    borderWidth: 2,
    borderColor: 'transparent',
  },
  inputContainerFocused: {
    borderColor: COLORS.primary,
    backgroundColor: COLORS.primaryMuted,
  },
  inputIcon: {
    width: 48,
    height: '100%',
    justifyContent: 'center',
    alignItems: 'center',
  },
  input: { 
    flex: 1,
    fontSize: TYPOGRAPHY.fontSize.md,
    color: COLORS.text,
    paddingRight: SPACING.base,
  },
  
  // Buttons
  nextButton: {
    backgroundColor: COLORS.primary,
    borderRadius: RADIUS.lg,
    height: 56,
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: SPACING.md,
    gap: SPACING.sm,
    ...SHADOWS.primary,
  },
  nextButtonText: {
    color: COLORS.secondary,
    fontSize: TYPOGRAPHY.fontSize.lg,
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  buttonRow: {
    flexDirection: 'row',
    gap: SPACING.md,
    marginTop: SPACING.md,
  },
  backStepButton: {
    flex: 0.4,
    backgroundColor: COLORS.background,
    borderRadius: RADIUS.lg,
    height: 56,
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    gap: SPACING.sm,
  },
  backStepButtonText: {
    color: COLORS.text,
    fontSize: TYPOGRAPHY.fontSize.md,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
  },
  registerButton: { 
    flex: 0.6,
    backgroundColor: COLORS.primary, 
    borderRadius: RADIUS.lg, 
    height: 56, 
    flexDirection: 'row',
    justifyContent: 'center', 
    alignItems: 'center',
    gap: SPACING.sm,
    ...SHADOWS.primary,
  },
  buttonDisabled: { 
    opacity: 0.5,
  },
  registerButtonText: { 
    color: COLORS.secondary, 
    fontSize: TYPOGRAPHY.fontSize.lg, 
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  
  // Login Link
  loginContainer: { 
    flexDirection: 'row', 
    justifyContent: 'center', 
    marginBottom: SPACING['2xl'],
  },
  loginText: { 
    color: COLORS.textSecondary,
    fontSize: TYPOGRAPHY.fontSize.base,
  },
  loginLink: { 
    color: COLORS.primary, 
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    fontSize: TYPOGRAPHY.fontSize.base,
  },
});
