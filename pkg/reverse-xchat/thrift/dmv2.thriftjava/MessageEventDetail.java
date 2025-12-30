package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.NoWhenBranchMatchedException;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.DefaultConstructorMarker;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000T\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0012\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\b6\u0018\u0000 \t2\u00020\u0001:\u0010\n\u000b\f\r\u000e\u000f\u0010\u0011\u0012\u0013\u0014\u0015\u0016\u0017\u0018\tB\t\b\u0004¢\u0006\u0004\b\u0002\u0010\u0003J\u0017\u0010\u0007\u001a\u00020\u00062\u0006\u0010\u0005\u001a\u00020\u0004H\u0016¢\u0006\u0004\b\u0007\u0010\b\u0082\u0001\u000e\u0019\u001a\u001b\u001c\u001d\u001e\u001f !\"#$%&¨\u0006'"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail;", "Lcom/bendb/thrifty/a;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "Companion", "MessageCreateEvent", "ConversationKeyChangeEvent", "GroupChangeEvent", "MessageFailureEvent", "MessageTypingEvent", "MessageDeleteEvent", "ConversationDeleteEvent", "ConversationMetadataChangeEvent", "GrokSearchResponseEvent", "RequestForEncryptedResendEvent", "MarkConversationReadEvent", "MarkConversationUnreadEvent", "MemberAccountDeleteEvent", "Unknown", "MessageEventDetailAdapter", "Lcom/x/dmv2/thriftjava/MessageEventDetail$ConversationDeleteEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$ConversationKeyChangeEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$ConversationMetadataChangeEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$GrokSearchResponseEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$GroupChangeEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$MarkConversationReadEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$MarkConversationUnreadEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$MemberAccountDeleteEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$MessageCreateEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$MessageDeleteEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$MessageFailureEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$MessageTypingEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$RequestForEncryptedResendEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail$Unknown;", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public abstract class MessageEventDetail implements InterfaceC11261a {

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MessageEventDetailAdapter();

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$ConversationDeleteEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/ConversationDeleteEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/ConversationDeleteEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/ConversationDeleteEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class ConversationDeleteEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.ConversationDeleteEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public ConversationDeleteEvent(@InterfaceC88464a com.x.dmv2.thriftjava.ConversationDeleteEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ ConversationDeleteEvent copy$default(ConversationDeleteEvent conversationDeleteEvent, com.x.dmv2.thriftjava.ConversationDeleteEvent conversationDeleteEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                conversationDeleteEvent2 = conversationDeleteEvent.value;
            }
            return conversationDeleteEvent.copy(conversationDeleteEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.ConversationDeleteEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final ConversationDeleteEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.ConversationDeleteEvent value) {
            Intrinsics.m65272h(value, "value");
            return new ConversationDeleteEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof ConversationDeleteEvent) && Intrinsics.m65267c(this.value, ((ConversationDeleteEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.ConversationDeleteEvent m76784getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(conversationDeleteEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$ConversationKeyChangeEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/ConversationKeyChangeEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/ConversationKeyChangeEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/ConversationKeyChangeEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class ConversationKeyChangeEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.ConversationKeyChangeEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public ConversationKeyChangeEvent(@InterfaceC88464a com.x.dmv2.thriftjava.ConversationKeyChangeEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ ConversationKeyChangeEvent copy$default(ConversationKeyChangeEvent conversationKeyChangeEvent, com.x.dmv2.thriftjava.ConversationKeyChangeEvent conversationKeyChangeEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                conversationKeyChangeEvent2 = conversationKeyChangeEvent.value;
            }
            return conversationKeyChangeEvent.copy(conversationKeyChangeEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.ConversationKeyChangeEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final ConversationKeyChangeEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.ConversationKeyChangeEvent value) {
            Intrinsics.m65272h(value, "value");
            return new ConversationKeyChangeEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof ConversationKeyChangeEvent) && Intrinsics.m65267c(this.value, ((ConversationKeyChangeEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.ConversationKeyChangeEvent m76785getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(conversationKeyChangeEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$ConversationMetadataChangeEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/ConversationMetadataChangeEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/ConversationMetadataChangeEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/ConversationMetadataChangeEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class ConversationMetadataChangeEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.ConversationMetadataChangeEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public ConversationMetadataChangeEvent(@InterfaceC88464a com.x.dmv2.thriftjava.ConversationMetadataChangeEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ ConversationMetadataChangeEvent copy$default(ConversationMetadataChangeEvent conversationMetadataChangeEvent, com.x.dmv2.thriftjava.ConversationMetadataChangeEvent conversationMetadataChangeEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                conversationMetadataChangeEvent2 = conversationMetadataChangeEvent.value;
            }
            return conversationMetadataChangeEvent.copy(conversationMetadataChangeEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.ConversationMetadataChangeEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final ConversationMetadataChangeEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.ConversationMetadataChangeEvent value) {
            Intrinsics.m65272h(value, "value");
            return new ConversationMetadataChangeEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof ConversationMetadataChangeEvent) && Intrinsics.m65267c(this.value, ((ConversationMetadataChangeEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.ConversationMetadataChangeEvent m76786getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(conversationMetadataChangeEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$GrokSearchResponseEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/GrokSearchResponseEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/GrokSearchResponseEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GrokSearchResponseEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GrokSearchResponseEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.GrokSearchResponseEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GrokSearchResponseEvent(@InterfaceC88464a com.x.dmv2.thriftjava.GrokSearchResponseEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GrokSearchResponseEvent copy$default(GrokSearchResponseEvent grokSearchResponseEvent, com.x.dmv2.thriftjava.GrokSearchResponseEvent grokSearchResponseEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                grokSearchResponseEvent2 = grokSearchResponseEvent.value;
            }
            return grokSearchResponseEvent.copy(grokSearchResponseEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.GrokSearchResponseEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GrokSearchResponseEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.GrokSearchResponseEvent value) {
            Intrinsics.m65272h(value, "value");
            return new GrokSearchResponseEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GrokSearchResponseEvent) && Intrinsics.m65267c(this.value, ((GrokSearchResponseEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.GrokSearchResponseEvent m76787getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(grokSearchResponseEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$GroupChangeEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/GroupChangeEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/GroupChangeEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/GroupChangeEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class GroupChangeEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.GroupChangeEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public GroupChangeEvent(@InterfaceC88464a com.x.dmv2.thriftjava.GroupChangeEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ GroupChangeEvent copy$default(GroupChangeEvent groupChangeEvent, com.x.dmv2.thriftjava.GroupChangeEvent groupChangeEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                groupChangeEvent2 = groupChangeEvent.value;
            }
            return groupChangeEvent.copy(groupChangeEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.GroupChangeEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final GroupChangeEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.GroupChangeEvent value) {
            Intrinsics.m65272h(value, "value");
            return new GroupChangeEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof GroupChangeEvent) && Intrinsics.m65267c(this.value, ((GroupChangeEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.GroupChangeEvent m76788getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(groupChangeEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$MarkConversationReadEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/MarkConversationReadEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/MarkConversationReadEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MarkConversationReadEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MarkConversationReadEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MarkConversationReadEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MarkConversationReadEvent(@InterfaceC88464a com.x.dmv2.thriftjava.MarkConversationReadEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MarkConversationReadEvent copy$default(MarkConversationReadEvent markConversationReadEvent, com.x.dmv2.thriftjava.MarkConversationReadEvent markConversationReadEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                markConversationReadEvent2 = markConversationReadEvent.value;
            }
            return markConversationReadEvent.copy(markConversationReadEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MarkConversationReadEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MarkConversationReadEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.MarkConversationReadEvent value) {
            Intrinsics.m65272h(value, "value");
            return new MarkConversationReadEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MarkConversationReadEvent) && Intrinsics.m65267c(this.value, ((MarkConversationReadEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MarkConversationReadEvent m76789getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(markConversationReadEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$MarkConversationUnreadEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/MarkConversationUnreadEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/MarkConversationUnreadEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MarkConversationUnreadEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MarkConversationUnreadEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MarkConversationUnreadEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MarkConversationUnreadEvent(@InterfaceC88464a com.x.dmv2.thriftjava.MarkConversationUnreadEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MarkConversationUnreadEvent copy$default(MarkConversationUnreadEvent markConversationUnreadEvent, com.x.dmv2.thriftjava.MarkConversationUnreadEvent markConversationUnreadEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                markConversationUnreadEvent2 = markConversationUnreadEvent.value;
            }
            return markConversationUnreadEvent.copy(markConversationUnreadEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MarkConversationUnreadEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MarkConversationUnreadEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.MarkConversationUnreadEvent value) {
            Intrinsics.m65272h(value, "value");
            return new MarkConversationUnreadEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MarkConversationUnreadEvent) && Intrinsics.m65267c(this.value, ((MarkConversationUnreadEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MarkConversationUnreadEvent m76790getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(markConversationUnreadEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$MemberAccountDeleteEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/MemberAccountDeleteEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/MemberAccountDeleteEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MemberAccountDeleteEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MemberAccountDeleteEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MemberAccountDeleteEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MemberAccountDeleteEvent(@InterfaceC88464a com.x.dmv2.thriftjava.MemberAccountDeleteEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MemberAccountDeleteEvent copy$default(MemberAccountDeleteEvent memberAccountDeleteEvent, com.x.dmv2.thriftjava.MemberAccountDeleteEvent memberAccountDeleteEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                memberAccountDeleteEvent2 = memberAccountDeleteEvent.value;
            }
            return memberAccountDeleteEvent.copy(memberAccountDeleteEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MemberAccountDeleteEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MemberAccountDeleteEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.MemberAccountDeleteEvent value) {
            Intrinsics.m65272h(value, "value");
            return new MemberAccountDeleteEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MemberAccountDeleteEvent) && Intrinsics.m65267c(this.value, ((MemberAccountDeleteEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MemberAccountDeleteEvent m76791getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(memberAccountDeleteEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$MessageCreateEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/MessageCreateEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageCreateEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageCreateEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MessageCreateEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MessageCreateEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MessageCreateEvent(@InterfaceC88464a com.x.dmv2.thriftjava.MessageCreateEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MessageCreateEvent copy$default(MessageCreateEvent messageCreateEvent, com.x.dmv2.thriftjava.MessageCreateEvent messageCreateEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                messageCreateEvent2 = messageCreateEvent.value;
            }
            return messageCreateEvent.copy(messageCreateEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MessageCreateEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MessageCreateEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.MessageCreateEvent value) {
            Intrinsics.m65272h(value, "value");
            return new MessageCreateEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MessageCreateEvent) && Intrinsics.m65267c(this.value, ((MessageCreateEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MessageCreateEvent m76792getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(messageCreateEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$MessageDeleteEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/MessageDeleteEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageDeleteEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageDeleteEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MessageDeleteEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MessageDeleteEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MessageDeleteEvent(@InterfaceC88464a com.x.dmv2.thriftjava.MessageDeleteEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MessageDeleteEvent copy$default(MessageDeleteEvent messageDeleteEvent, com.x.dmv2.thriftjava.MessageDeleteEvent messageDeleteEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                messageDeleteEvent2 = messageDeleteEvent.value;
            }
            return messageDeleteEvent.copy(messageDeleteEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MessageDeleteEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MessageDeleteEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.MessageDeleteEvent value) {
            Intrinsics.m65272h(value, "value");
            return new MessageDeleteEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MessageDeleteEvent) && Intrinsics.m65267c(this.value, ((MessageDeleteEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MessageDeleteEvent m76793getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(messageDeleteEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$MessageEventDetailAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MessageEventDetail;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MessageEventDetail;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MessageEventDetailAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MessageEventDetail m85659read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            MessageEventDetail messageCreateEvent;
            Intrinsics.m65272h(protocol, "protocol");
            MessageEventDetail messageEventDetail = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    if (messageEventDetail != null) {
                        return messageEventDetail;
                    }
                    throw new IllegalStateException("unreadable");
                }
                switch (c11265cMo14127V2.f38393b) {
                    case 1:
                        if (b == 12) {
                            messageCreateEvent = new MessageCreateEvent((com.x.dmv2.thriftjava.MessageCreateEvent) com.x.dmv2.thriftjava.MessageCreateEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 2:
                    default:
                        messageEventDetail = Unknown.INSTANCE;
                        C11272a.m14141a(protocol, b);
                        continue;
                    case 3:
                        if (b == 12) {
                            messageCreateEvent = new ConversationKeyChangeEvent((com.x.dmv2.thriftjava.ConversationKeyChangeEvent) com.x.dmv2.thriftjava.ConversationKeyChangeEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 4:
                        if (b == 12) {
                            messageCreateEvent = new GroupChangeEvent((com.x.dmv2.thriftjava.GroupChangeEvent) com.x.dmv2.thriftjava.GroupChangeEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 5:
                        if (b == 12) {
                            messageCreateEvent = new MessageFailureEvent((com.x.dmv2.thriftjava.MessageFailureEvent) com.x.dmv2.thriftjava.MessageFailureEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 6:
                        if (b == 12) {
                            messageCreateEvent = new MessageTypingEvent((com.x.dmv2.thriftjava.MessageTypingEvent) com.x.dmv2.thriftjava.MessageTypingEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 7:
                        if (b == 12) {
                            messageCreateEvent = new MessageDeleteEvent((com.x.dmv2.thriftjava.MessageDeleteEvent) com.x.dmv2.thriftjava.MessageDeleteEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 8:
                        if (b == 12) {
                            messageCreateEvent = new ConversationDeleteEvent((com.x.dmv2.thriftjava.ConversationDeleteEvent) com.x.dmv2.thriftjava.ConversationDeleteEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 9:
                        if (b == 12) {
                            messageCreateEvent = new ConversationMetadataChangeEvent((com.x.dmv2.thriftjava.ConversationMetadataChangeEvent) com.x.dmv2.thriftjava.ConversationMetadataChangeEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 10:
                        if (b == 12) {
                            messageCreateEvent = new GrokSearchResponseEvent((com.x.dmv2.thriftjava.GrokSearchResponseEvent) com.x.dmv2.thriftjava.GrokSearchResponseEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 11:
                        if (b == 12) {
                            messageCreateEvent = new RequestForEncryptedResendEvent((com.x.dmv2.thriftjava.RequestForEncryptedResendEvent) com.x.dmv2.thriftjava.RequestForEncryptedResendEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 12:
                        if (b == 12) {
                            messageCreateEvent = new MarkConversationReadEvent((com.x.dmv2.thriftjava.MarkConversationReadEvent) com.x.dmv2.thriftjava.MarkConversationReadEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 13:
                        if (b == 12) {
                            messageCreateEvent = new MarkConversationUnreadEvent((com.x.dmv2.thriftjava.MarkConversationUnreadEvent) com.x.dmv2.thriftjava.MarkConversationUnreadEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 14:
                        if (b == 12) {
                            messageCreateEvent = new MemberAccountDeleteEvent((com.x.dmv2.thriftjava.MemberAccountDeleteEvent) com.x.dmv2.thriftjava.MemberAccountDeleteEvent.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                }
                messageEventDetail = messageCreateEvent;
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MessageEventDetail struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MessageEventDetail");
            if (struct instanceof MessageCreateEvent) {
                protocol.mo14136v3("messageCreateEvent", 1, (byte) 12);
                com.x.dmv2.thriftjava.MessageCreateEvent.ADAPTER.write(protocol, ((MessageCreateEvent) struct).m76792getValue());
            } else if (struct instanceof ConversationKeyChangeEvent) {
                protocol.mo14136v3("conversationKeyChangeEvent", 3, (byte) 12);
                com.x.dmv2.thriftjava.ConversationKeyChangeEvent.ADAPTER.write(protocol, ((ConversationKeyChangeEvent) struct).m76785getValue());
            } else if (struct instanceof GroupChangeEvent) {
                protocol.mo14136v3("groupChangeEvent", 4, (byte) 12);
                com.x.dmv2.thriftjava.GroupChangeEvent.ADAPTER.write(protocol, ((GroupChangeEvent) struct).m76788getValue());
            } else if (struct instanceof MessageFailureEvent) {
                protocol.mo14136v3("messageFailureEvent", 5, (byte) 12);
                com.x.dmv2.thriftjava.MessageFailureEvent.ADAPTER.write(protocol, ((MessageFailureEvent) struct).m76794getValue());
            } else if (struct instanceof MessageTypingEvent) {
                protocol.mo14136v3("messageTypingEvent", 6, (byte) 12);
                com.x.dmv2.thriftjava.MessageTypingEvent.ADAPTER.write(protocol, ((MessageTypingEvent) struct).m76795getValue());
            } else if (struct instanceof MessageDeleteEvent) {
                protocol.mo14136v3("messageDeleteEvent", 7, (byte) 12);
                com.x.dmv2.thriftjava.MessageDeleteEvent.ADAPTER.write(protocol, ((MessageDeleteEvent) struct).m76793getValue());
            } else if (struct instanceof ConversationDeleteEvent) {
                protocol.mo14136v3("conversationDeleteEvent", 8, (byte) 12);
                com.x.dmv2.thriftjava.ConversationDeleteEvent.ADAPTER.write(protocol, ((ConversationDeleteEvent) struct).m76784getValue());
            } else if (struct instanceof ConversationMetadataChangeEvent) {
                protocol.mo14136v3("conversationMetadataChangeEvent", 9, (byte) 12);
                com.x.dmv2.thriftjava.ConversationMetadataChangeEvent.ADAPTER.write(protocol, ((ConversationMetadataChangeEvent) struct).m76786getValue());
            } else if (struct instanceof GrokSearchResponseEvent) {
                protocol.mo14136v3("grokSearchResponseEvent", 10, (byte) 12);
                com.x.dmv2.thriftjava.GrokSearchResponseEvent.ADAPTER.write(protocol, ((GrokSearchResponseEvent) struct).m76787getValue());
            } else if (struct instanceof RequestForEncryptedResendEvent) {
                protocol.mo14136v3("requestForEncryptedResendEvent", 11, (byte) 12);
                com.x.dmv2.thriftjava.RequestForEncryptedResendEvent.ADAPTER.write(protocol, ((RequestForEncryptedResendEvent) struct).m76796getValue());
            } else if (struct instanceof MarkConversationReadEvent) {
                protocol.mo14136v3("markConversationReadEvent", 12, (byte) 12);
                com.x.dmv2.thriftjava.MarkConversationReadEvent.ADAPTER.write(protocol, ((MarkConversationReadEvent) struct).m76789getValue());
            } else if (struct instanceof MarkConversationUnreadEvent) {
                protocol.mo14136v3("markConversationUnreadEvent", 13, (byte) 12);
                com.x.dmv2.thriftjava.MarkConversationUnreadEvent.ADAPTER.write(protocol, ((MarkConversationUnreadEvent) struct).m76790getValue());
            } else if (struct instanceof MemberAccountDeleteEvent) {
                protocol.mo14136v3("memberAccountDeleteEvent", 14, (byte) 12);
                com.x.dmv2.thriftjava.MemberAccountDeleteEvent.ADAPTER.write(protocol, ((MemberAccountDeleteEvent) struct).m76791getValue());
            } else if (!(struct instanceof Unknown)) {
                throw new NoWhenBranchMatchedException();
            }
            protocol.mo14134i0();
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$MessageFailureEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/MessageFailureEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageFailureEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageFailureEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MessageFailureEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MessageFailureEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MessageFailureEvent(@InterfaceC88464a com.x.dmv2.thriftjava.MessageFailureEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MessageFailureEvent copy$default(MessageFailureEvent messageFailureEvent, com.x.dmv2.thriftjava.MessageFailureEvent messageFailureEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                messageFailureEvent2 = messageFailureEvent.value;
            }
            return messageFailureEvent.copy(messageFailureEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MessageFailureEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MessageFailureEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.MessageFailureEvent value) {
            Intrinsics.m65272h(value, "value");
            return new MessageFailureEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MessageFailureEvent) && Intrinsics.m65267c(this.value, ((MessageFailureEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MessageFailureEvent m76794getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(messageFailureEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$MessageTypingEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/MessageTypingEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageTypingEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageTypingEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MessageTypingEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MessageTypingEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MessageTypingEvent(@InterfaceC88464a com.x.dmv2.thriftjava.MessageTypingEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MessageTypingEvent copy$default(MessageTypingEvent messageTypingEvent, com.x.dmv2.thriftjava.MessageTypingEvent messageTypingEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                messageTypingEvent2 = messageTypingEvent.value;
            }
            return messageTypingEvent.copy(messageTypingEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MessageTypingEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MessageTypingEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.MessageTypingEvent value) {
            Intrinsics.m65272h(value, "value");
            return new MessageTypingEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MessageTypingEvent) && Intrinsics.m65267c(this.value, ((MessageTypingEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MessageTypingEvent m76795getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(messageTypingEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$RequestForEncryptedResendEvent;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "value", "Lcom/x/dmv2/thriftjava/RequestForEncryptedResendEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/RequestForEncryptedResendEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/RequestForEncryptedResendEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class RequestForEncryptedResendEvent extends MessageEventDetail {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.RequestForEncryptedResendEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public RequestForEncryptedResendEvent(@InterfaceC88464a com.x.dmv2.thriftjava.RequestForEncryptedResendEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ RequestForEncryptedResendEvent copy$default(RequestForEncryptedResendEvent requestForEncryptedResendEvent, com.x.dmv2.thriftjava.RequestForEncryptedResendEvent requestForEncryptedResendEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                requestForEncryptedResendEvent2 = requestForEncryptedResendEvent.value;
            }
            return requestForEncryptedResendEvent.copy(requestForEncryptedResendEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.RequestForEncryptedResendEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final RequestForEncryptedResendEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.RequestForEncryptedResendEvent value) {
            Intrinsics.m65272h(value, "value");
            return new RequestForEncryptedResendEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof RequestForEncryptedResendEvent) && Intrinsics.m65267c(this.value, ((RequestForEncryptedResendEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.RequestForEncryptedResendEvent m76796getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEventDetail(requestForEncryptedResendEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\n\u0002\u0010\u000e\n\u0000\bÆ\n\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0013\u0010\u0004\u001a\u00020\u00052\b\u0010\u0006\u001a\u0004\u0018\u00010\u0007HÖ\u0003J\t\u0010\b\u001a\u00020\tHÖ\u0001J\t\u0010\n\u001a\u00020\u000bHÖ\u0001¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEventDetail$Unknown;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "<init>", "()V", "equals", "", "other", "", "hashCode", "", "toString", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Unknown extends MessageEventDetail {

        @InterfaceC88464a
        public static final Unknown INSTANCE = new Unknown();

        private Unknown() {
            super(null);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            return this == other || (other instanceof Unknown);
        }

        public int hashCode() {
            return -2013842451;
        }

        @InterfaceC88464a
        public String toString() {
            return "Unknown";
        }
    }

    public /* synthetic */ MessageEventDetail(DefaultConstructorMarker defaultConstructorMarker) {
        this();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }

    private MessageEventDetail() {
    }
}
