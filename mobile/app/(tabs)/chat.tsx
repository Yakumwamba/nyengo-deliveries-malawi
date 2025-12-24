import React from 'react';
import { View, Text, StyleSheet, FlatList, TouchableOpacity } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { COLORS, TYPOGRAPHY, SPACING, RADIUS, SHADOWS } from '../../src/constants/theme';

const mockChats = [
  { id: '1', name: 'John Doe', lastMessage: 'I am almost there!', time: '2m ago', unread: 2, orderNumber: 'NYG-ABC123', status: 'active' },
  { id: '2', name: 'Jane Smith', lastMessage: 'Thank you for the delivery', time: '1h ago', unread: 0, orderNumber: 'NYG-DEF456', status: 'completed' },
  { id: '3', name: 'Bob Wilson', lastMessage: 'Can you call me when you arrive?', time: '3h ago', unread: 1, orderNumber: 'NYG-GHI789', status: 'active' },
  { id: '4', name: 'Alice Brown', lastMessage: 'Package received, thanks!', time: 'Yesterday', unread: 0, orderNumber: 'NYG-JKL012', status: 'completed' },
];

export default function ChatScreen() {
  return (
    <View style={styles.container}>
      {/* Header Stats */}
      <View style={styles.headerStats}>
        <View style={styles.statItem}>
          <Text style={styles.statValue}>4</Text>
          <Text style={styles.statLabel}>Active Chats</Text>
        </View>
        <View style={styles.statDivider} />
        <View style={styles.statItem}>
          <Text style={styles.statValue}>3</Text>
          <Text style={styles.statLabel}>Unread</Text>
        </View>
      </View>

      <FlatList
        data={mockChats}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => <ChatItem chat={item} />}
        ListEmptyComponent={
          <View style={styles.emptyContainer}>
            <View style={styles.emptyIcon}>
              <Ionicons name="chatbubbles-outline" size={56} color={COLORS.textMuted} />
            </View>
            <Text style={styles.emptyTitle}>No conversations</Text>
            <Text style={styles.emptyText}>
              Start a delivery to chat with customers
            </Text>
          </View>
        }
        contentContainerStyle={styles.list}
        showsVerticalScrollIndicator={false}
      />
    </View>
  );
}

function ChatItem({ chat }: { chat: typeof mockChats[0] }) {
  const isActive = chat.status === 'active';
  
  return (
    <TouchableOpacity style={styles.chatItem} activeOpacity={0.7}>
      <View style={styles.avatarContainer}>
        <View style={[styles.avatar, isActive && styles.avatarActive]}>
          <Text style={styles.avatarText}>{chat.name.charAt(0)}</Text>
        </View>
        {isActive && <View style={styles.onlineIndicator} />}
      </View>
      
      <View style={styles.chatContent}>
        <View style={styles.chatHeader}>
          <Text style={styles.chatName}>{chat.name}</Text>
          <Text style={[styles.chatTime, chat.unread > 0 && styles.chatTimeUnread]}>
            {chat.time}
          </Text>
        </View>
        
        <View style={styles.orderBadge}>
          <Ionicons name="cube-outline" size={12} color={COLORS.primary} />
          <Text style={styles.orderNumber}>{chat.orderNumber}</Text>
        </View>
        
        <View style={styles.chatFooter}>
          <Text style={styles.lastMessage} numberOfLines={1}>
            {chat.lastMessage}
          </Text>
          {chat.unread > 0 && (
            <View style={styles.unreadBadge}>
              <Text style={styles.unreadText}>{chat.unread}</Text>
            </View>
          )}
        </View>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  container: { 
    flex: 1, 
    backgroundColor: COLORS.background,
  },
  
  // Header Stats
  headerStats: {
    flexDirection: 'row',
    backgroundColor: COLORS.surface,
    marginHorizontal: SPACING.base,
    marginTop: SPACING.base,
    padding: SPACING.base,
    borderRadius: RADIUS.lg,
    ...SHADOWS.sm,
  },
  statItem: {
    flex: 1,
    alignItems: 'center',
  },
  statValue: {
    fontSize: TYPOGRAPHY.fontSize['2xl'],
    fontWeight: TYPOGRAPHY.fontWeight.bold,
    color: COLORS.text,
  },
  statLabel: {
    fontSize: TYPOGRAPHY.fontSize.xs,
    color: COLORS.textMuted,
    marginTop: 2,
  },
  statDivider: {
    width: 1,
    backgroundColor: COLORS.border,
  },
  
  // List
  list: { 
    padding: SPACING.base,
    paddingTop: SPACING.md,
  },
  
  // Chat Item
  chatItem: { 
    backgroundColor: COLORS.surface, 
    borderRadius: RADIUS.xl, 
    padding: SPACING.base, 
    marginBottom: SPACING.md,
    flexDirection: 'row', 
    alignItems: 'flex-start',
    ...SHADOWS.sm,
  },
  avatarContainer: {
    position: 'relative',
    marginRight: SPACING.md,
  },
  avatar: { 
    width: 52, 
    height: 52, 
    borderRadius: 26, 
    backgroundColor: COLORS.secondaryLight, 
    justifyContent: 'center', 
    alignItems: 'center',
  },
  avatarActive: {
    backgroundColor: COLORS.primary,
  },
  avatarText: { 
    color: COLORS.white, 
    fontSize: TYPOGRAPHY.fontSize.xl, 
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  onlineIndicator: {
    position: 'absolute',
    bottom: 2,
    right: 2,
    width: 14,
    height: 14,
    borderRadius: 7,
    backgroundColor: COLORS.success,
    borderWidth: 2,
    borderColor: COLORS.surface,
  },
  chatContent: { 
    flex: 1,
  },
  chatHeader: { 
    flexDirection: 'row', 
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 4,
  },
  chatName: { 
    fontSize: TYPOGRAPHY.fontSize.md, 
    fontWeight: TYPOGRAPHY.fontWeight.semibold, 
    color: COLORS.text,
  },
  chatTime: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    color: COLORS.textMuted,
  },
  chatTimeUnread: {
    color: COLORS.primary,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
  },
  orderBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    alignSelf: 'flex-start',
    backgroundColor: COLORS.primary + '15',
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: RADIUS.sm,
    marginBottom: 6,
    gap: 4,
  },
  orderNumber: { 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    color: COLORS.primary,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
  },
  chatFooter: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  lastMessage: { 
    fontSize: TYPOGRAPHY.fontSize.base, 
    color: COLORS.textSecondary,
    flex: 1,
    marginRight: SPACING.sm,
  },
  unreadBadge: { 
    backgroundColor: COLORS.accent, 
    minWidth: 22, 
    height: 22, 
    borderRadius: 11, 
    justifyContent: 'center', 
    alignItems: 'center',
    paddingHorizontal: 6,
  },
  unreadText: { 
    color: COLORS.white, 
    fontSize: TYPOGRAPHY.fontSize.xs, 
    fontWeight: TYPOGRAPHY.fontWeight.bold,
  },
  
  // Empty State
  emptyContainer: { 
    alignItems: 'center', 
    paddingTop: SPACING['4xl'],
  },
  emptyIcon: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: COLORS.surface,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: SPACING.lg,
    ...SHADOWS.md,
  },
  emptyTitle: {
    fontSize: TYPOGRAPHY.fontSize.lg,
    fontWeight: TYPOGRAPHY.fontWeight.semibold,
    color: COLORS.text,
    marginBottom: SPACING.sm,
  },
  emptyText: { 
    fontSize: TYPOGRAPHY.fontSize.base, 
    color: COLORS.textMuted,
    textAlign: 'center',
    paddingHorizontal: SPACING['2xl'],
  },
});
