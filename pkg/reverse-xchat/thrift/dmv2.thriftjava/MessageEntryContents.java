package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import com.socure.docv.capturesdk.common.utils.ApiConstant;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.NoWhenBranchMatchedException;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.DefaultConstructorMarker;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000`\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0015\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\b6\u0018\u0000 \t2\u00020\u0001:\u0013\n\u000b\f\r\u000e\u000f\u0010\u0011\u0012\u0013\u0014\u0015\u0016\u0017\u0018\u0019\u001a\u001b\tB\t\b\u0004¢\u0006\u0004\b\u0002\u0010\u0003J\u0017\u0010\u0007\u001a\u00020\u00062\u0006\u0010\u0005\u001a\u00020\u0004H\u0016¢\u0006\u0004\b\u0007\u0010\b\u0082\u0001\u0011\u001c\u001d\u001e\u001f !\"#$%&'()*+,¨\u0006-"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents;", "Lcom/bendb/thrifty/a;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "Companion", "Message", "ReactionAdd", "ReactionRemove", "MessageEdit", "MarkConversationRead", "MarkConversationUnread", "PinConversation", "UnpinConversation", "ScreenCaptureDetected", "AvCallEnded", "AvCallMissed", "DraftMessage", "AcceptMessageRequest", "NicknameMessage", "SetVerifiedStatus", "AvCallStarted", "Unknown", "MessageEntryContentsAdapter", "Lcom/x/dmv2/thriftjava/MessageEntryContents$AcceptMessageRequest;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$AvCallEnded;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$AvCallMissed;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$AvCallStarted;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$DraftMessage;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$MarkConversationRead;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$MarkConversationUnread;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$Message;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$MessageEdit;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$NicknameMessage;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$PinConversation;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$ReactionAdd;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$ReactionRemove;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$ScreenCaptureDetected;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$SetVerifiedStatus;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$Unknown;", "Lcom/x/dmv2/thriftjava/MessageEntryContents$UnpinConversation;", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public abstract class MessageEntryContents implements InterfaceC11261a {

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MessageEntryContentsAdapter();

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$AcceptMessageRequest;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/AcceptMessageRequest;", "<init>", "(Lcom/x/dmv2/thriftjava/AcceptMessageRequest;)V", "getValue", "()Lcom/x/dmv2/thriftjava/AcceptMessageRequest;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class AcceptMessageRequest extends MessageEntryContents {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.AcceptMessageRequest value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public AcceptMessageRequest(@InterfaceC88464a com.x.dmv2.thriftjava.AcceptMessageRequest value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ AcceptMessageRequest copy$default(AcceptMessageRequest acceptMessageRequest, com.x.dmv2.thriftjava.AcceptMessageRequest acceptMessageRequest2, int i, Object obj) {
            if ((i & 1) != 0) {
                acceptMessageRequest2 = acceptMessageRequest.value;
            }
            return acceptMessageRequest.copy(acceptMessageRequest2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.AcceptMessageRequest getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final AcceptMessageRequest copy(@InterfaceC88464a com.x.dmv2.thriftjava.AcceptMessageRequest value) {
            Intrinsics.m65272h(value, "value");
            return new AcceptMessageRequest(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof AcceptMessageRequest) && Intrinsics.m65267c(this.value, ((AcceptMessageRequest) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.AcceptMessageRequest m76768getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(accept_message_request=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$AvCallEnded;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/AVCallEnded;", "<init>", "(Lcom/x/dmv2/thriftjava/AVCallEnded;)V", "getValue", "()Lcom/x/dmv2/thriftjava/AVCallEnded;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class AvCallEnded extends MessageEntryContents {

        @InterfaceC88464a
        private final AVCallEnded value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public AvCallEnded(@InterfaceC88464a AVCallEnded value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ AvCallEnded copy$default(AvCallEnded avCallEnded, AVCallEnded aVCallEnded, int i, Object obj) {
            if ((i & 1) != 0) {
                aVCallEnded = avCallEnded.value;
            }
            return avCallEnded.copy(aVCallEnded);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final AVCallEnded getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final AvCallEnded copy(@InterfaceC88464a AVCallEnded value) {
            Intrinsics.m65272h(value, "value");
            return new AvCallEnded(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof AvCallEnded) && Intrinsics.m65267c(this.value, ((AvCallEnded) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final AVCallEnded m76769getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(av_call_ended=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$AvCallMissed;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/AVCallMissed;", "<init>", "(Lcom/x/dmv2/thriftjava/AVCallMissed;)V", "getValue", "()Lcom/x/dmv2/thriftjava/AVCallMissed;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class AvCallMissed extends MessageEntryContents {

        @InterfaceC88464a
        private final AVCallMissed value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public AvCallMissed(@InterfaceC88464a AVCallMissed value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ AvCallMissed copy$default(AvCallMissed avCallMissed, AVCallMissed aVCallMissed, int i, Object obj) {
            if ((i & 1) != 0) {
                aVCallMissed = avCallMissed.value;
            }
            return avCallMissed.copy(aVCallMissed);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final AVCallMissed getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final AvCallMissed copy(@InterfaceC88464a AVCallMissed value) {
            Intrinsics.m65272h(value, "value");
            return new AvCallMissed(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof AvCallMissed) && Intrinsics.m65267c(this.value, ((AvCallMissed) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final AVCallMissed m76770getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(av_call_missed=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$AvCallStarted;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/AVCallStarted;", "<init>", "(Lcom/x/dmv2/thriftjava/AVCallStarted;)V", "getValue", "()Lcom/x/dmv2/thriftjava/AVCallStarted;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class AvCallStarted extends MessageEntryContents {

        @InterfaceC88464a
        private final AVCallStarted value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public AvCallStarted(@InterfaceC88464a AVCallStarted value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ AvCallStarted copy$default(AvCallStarted avCallStarted, AVCallStarted aVCallStarted, int i, Object obj) {
            if ((i & 1) != 0) {
                aVCallStarted = avCallStarted.value;
            }
            return avCallStarted.copy(aVCallStarted);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final AVCallStarted getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final AvCallStarted copy(@InterfaceC88464a AVCallStarted value) {
            Intrinsics.m65272h(value, "value");
            return new AvCallStarted(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof AvCallStarted) && Intrinsics.m65267c(this.value, ((AvCallStarted) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final AVCallStarted m76771getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(av_call_started=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$DraftMessage;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/DraftMessage;", "<init>", "(Lcom/x/dmv2/thriftjava/DraftMessage;)V", "getValue", "()Lcom/x/dmv2/thriftjava/DraftMessage;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class DraftMessage extends MessageEntryContents {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.DraftMessage value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public DraftMessage(@InterfaceC88464a com.x.dmv2.thriftjava.DraftMessage value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ DraftMessage copy$default(DraftMessage draftMessage, com.x.dmv2.thriftjava.DraftMessage draftMessage2, int i, Object obj) {
            if ((i & 1) != 0) {
                draftMessage2 = draftMessage.value;
            }
            return draftMessage.copy(draftMessage2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.DraftMessage getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final DraftMessage copy(@InterfaceC88464a com.x.dmv2.thriftjava.DraftMessage value) {
            Intrinsics.m65272h(value, "value");
            return new DraftMessage(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof DraftMessage) && Intrinsics.m65267c(this.value, ((DraftMessage) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.DraftMessage m76772getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(draft_message=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$MarkConversationRead;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/MarkConversationRead;", "<init>", "(Lcom/x/dmv2/thriftjava/MarkConversationRead;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MarkConversationRead;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MarkConversationRead extends MessageEntryContents {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MarkConversationRead value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MarkConversationRead(@InterfaceC88464a com.x.dmv2.thriftjava.MarkConversationRead value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MarkConversationRead copy$default(MarkConversationRead markConversationRead, com.x.dmv2.thriftjava.MarkConversationRead markConversationRead2, int i, Object obj) {
            if ((i & 1) != 0) {
                markConversationRead2 = markConversationRead.value;
            }
            return markConversationRead.copy(markConversationRead2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MarkConversationRead getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MarkConversationRead copy(@InterfaceC88464a com.x.dmv2.thriftjava.MarkConversationRead value) {
            Intrinsics.m65272h(value, "value");
            return new MarkConversationRead(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MarkConversationRead) && Intrinsics.m65267c(this.value, ((MarkConversationRead) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MarkConversationRead m76773getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(mark_conversation_read=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$MarkConversationUnread;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/MarkConversationUnread;", "<init>", "(Lcom/x/dmv2/thriftjava/MarkConversationUnread;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MarkConversationUnread;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MarkConversationUnread extends MessageEntryContents {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MarkConversationUnread value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MarkConversationUnread(@InterfaceC88464a com.x.dmv2.thriftjava.MarkConversationUnread value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MarkConversationUnread copy$default(MarkConversationUnread markConversationUnread, com.x.dmv2.thriftjava.MarkConversationUnread markConversationUnread2, int i, Object obj) {
            if ((i & 1) != 0) {
                markConversationUnread2 = markConversationUnread.value;
            }
            return markConversationUnread.copy(markConversationUnread2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MarkConversationUnread getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MarkConversationUnread copy(@InterfaceC88464a com.x.dmv2.thriftjava.MarkConversationUnread value) {
            Intrinsics.m65272h(value, "value");
            return new MarkConversationUnread(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MarkConversationUnread) && Intrinsics.m65267c(this.value, ((MarkConversationUnread) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MarkConversationUnread m76774getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(mark_conversation_unread=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$Message;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/MessageContents;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageContents;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageContents;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Message extends MessageEntryContents {

        @InterfaceC88464a
        private final MessageContents value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Message(@InterfaceC88464a MessageContents value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Message copy$default(Message message, MessageContents messageContents, int i, Object obj) {
            if ((i & 1) != 0) {
                messageContents = message.value;
            }
            return message.copy(messageContents);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final MessageContents getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Message copy(@InterfaceC88464a MessageContents value) {
            Intrinsics.m65272h(value, "value");
            return new Message(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Message) && Intrinsics.m65267c(this.value, ((Message) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final MessageContents m76775getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(message=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$MessageEdit;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/MessageEdit;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageEdit;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageEdit;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MessageEdit extends MessageEntryContents {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MessageEdit value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MessageEdit(@InterfaceC88464a com.x.dmv2.thriftjava.MessageEdit value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MessageEdit copy$default(MessageEdit messageEdit, com.x.dmv2.thriftjava.MessageEdit messageEdit2, int i, Object obj) {
            if ((i & 1) != 0) {
                messageEdit2 = messageEdit.value;
            }
            return messageEdit.copy(messageEdit2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MessageEdit getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MessageEdit copy(@InterfaceC88464a com.x.dmv2.thriftjava.MessageEdit value) {
            Intrinsics.m65272h(value, "value");
            return new MessageEdit(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MessageEdit) && Intrinsics.m65267c(this.value, ((MessageEdit) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MessageEdit m76776getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(message_edit=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$MessageEntryContentsAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MessageEntryContents;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MessageEntryContents;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MessageEntryContentsAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MessageEntryContents m85656read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            MessageEntryContents message;
            Intrinsics.m65272h(protocol, "protocol");
            MessageEntryContents messageEntryContents = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    if (messageEntryContents != null) {
                        return messageEntryContents;
                    }
                    throw new IllegalStateException("unreadable");
                }
                switch (c11265cMo14127V2.f38393b) {
                    case 1:
                        if (b == 12) {
                            message = new Message((MessageContents) MessageContents.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 2:
                        if (b == 12) {
                            message = new ReactionAdd((MessageReactionAdd) MessageReactionAdd.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 3:
                        if (b == 12) {
                            message = new ReactionRemove((MessageReactionRemove) MessageReactionRemove.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 4:
                        if (b == 12) {
                            message = new MessageEdit((com.x.dmv2.thriftjava.MessageEdit) com.x.dmv2.thriftjava.MessageEdit.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 5:
                        if (b == 12) {
                            message = new MarkConversationRead((com.x.dmv2.thriftjava.MarkConversationRead) com.x.dmv2.thriftjava.MarkConversationRead.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 6:
                        if (b == 12) {
                            message = new MarkConversationUnread((com.x.dmv2.thriftjava.MarkConversationUnread) com.x.dmv2.thriftjava.MarkConversationUnread.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 7:
                        if (b == 12) {
                            message = new PinConversation((com.x.dmv2.thriftjava.PinConversation) com.x.dmv2.thriftjava.PinConversation.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 8:
                        if (b == 12) {
                            message = new UnpinConversation((com.x.dmv2.thriftjava.UnpinConversation) com.x.dmv2.thriftjava.UnpinConversation.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 9:
                        if (b == 12) {
                            message = new ScreenCaptureDetected((com.x.dmv2.thriftjava.ScreenCaptureDetected) com.x.dmv2.thriftjava.ScreenCaptureDetected.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 10:
                        if (b == 12) {
                            message = new AvCallEnded((AVCallEnded) AVCallEnded.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 11:
                        if (b == 12) {
                            message = new AvCallMissed((AVCallMissed) AVCallMissed.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 12:
                        if (b == 12) {
                            message = new DraftMessage((com.x.dmv2.thriftjava.DraftMessage) com.x.dmv2.thriftjava.DraftMessage.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 13:
                        if (b == 12) {
                            message = new AcceptMessageRequest((com.x.dmv2.thriftjava.AcceptMessageRequest) com.x.dmv2.thriftjava.AcceptMessageRequest.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 14:
                        if (b == 12) {
                            message = new NicknameMessage((com.x.dmv2.thriftjava.NicknameMessage) com.x.dmv2.thriftjava.NicknameMessage.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 15:
                        if (b == 12) {
                            message = new SetVerifiedStatus((com.x.dmv2.thriftjava.SetVerifiedStatus) com.x.dmv2.thriftjava.SetVerifiedStatus.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 16:
                        if (b == 12) {
                            message = new AvCallStarted((AVCallStarted) AVCallStarted.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    default:
                        messageEntryContents = Unknown.INSTANCE;
                        C11272a.m14141a(protocol, b);
                        continue;
                }
                messageEntryContents = message;
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MessageEntryContents struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MessageEntryContents");
            if (struct instanceof Message) {
                protocol.mo14136v3(ApiConstant.KEY_MESSAGE, 1, (byte) 12);
                MessageContents.ADAPTER.write(protocol, ((Message) struct).m76775getValue());
            } else if (struct instanceof ReactionAdd) {
                protocol.mo14136v3("reaction_add", 2, (byte) 12);
                MessageReactionAdd.ADAPTER.write(protocol, ((ReactionAdd) struct).m76779getValue());
            } else if (struct instanceof ReactionRemove) {
                protocol.mo14136v3("reaction_remove", 3, (byte) 12);
                MessageReactionRemove.ADAPTER.write(protocol, ((ReactionRemove) struct).m76780getValue());
            } else if (struct instanceof MessageEdit) {
                protocol.mo14136v3("message_edit", 4, (byte) 12);
                com.x.dmv2.thriftjava.MessageEdit.ADAPTER.write(protocol, ((MessageEdit) struct).m76776getValue());
            } else if (struct instanceof MarkConversationRead) {
                protocol.mo14136v3("mark_conversation_read", 5, (byte) 12);
                com.x.dmv2.thriftjava.MarkConversationRead.ADAPTER.write(protocol, ((MarkConversationRead) struct).m76773getValue());
            } else if (struct instanceof MarkConversationUnread) {
                protocol.mo14136v3("mark_conversation_unread", 6, (byte) 12);
                com.x.dmv2.thriftjava.MarkConversationUnread.ADAPTER.write(protocol, ((MarkConversationUnread) struct).m76774getValue());
            } else if (struct instanceof PinConversation) {
                protocol.mo14136v3("pin_conversation", 7, (byte) 12);
                com.x.dmv2.thriftjava.PinConversation.ADAPTER.write(protocol, ((PinConversation) struct).m76778getValue());
            } else if (struct instanceof UnpinConversation) {
                protocol.mo14136v3("unpin_conversation", 8, (byte) 12);
                com.x.dmv2.thriftjava.UnpinConversation.ADAPTER.write(protocol, ((UnpinConversation) struct).m76783getValue());
            } else if (struct instanceof ScreenCaptureDetected) {
                protocol.mo14136v3("screen_capture_detected", 9, (byte) 12);
                com.x.dmv2.thriftjava.ScreenCaptureDetected.ADAPTER.write(protocol, ((ScreenCaptureDetected) struct).m76781getValue());
            } else if (struct instanceof AvCallEnded) {
                protocol.mo14136v3("av_call_ended", 10, (byte) 12);
                AVCallEnded.ADAPTER.write(protocol, ((AvCallEnded) struct).m76769getValue());
            } else if (struct instanceof AvCallMissed) {
                protocol.mo14136v3("av_call_missed", 11, (byte) 12);
                AVCallMissed.ADAPTER.write(protocol, ((AvCallMissed) struct).m76770getValue());
            } else if (struct instanceof DraftMessage) {
                protocol.mo14136v3("draft_message", 12, (byte) 12);
                com.x.dmv2.thriftjava.DraftMessage.ADAPTER.write(protocol, ((DraftMessage) struct).m76772getValue());
            } else if (struct instanceof AcceptMessageRequest) {
                protocol.mo14136v3("accept_message_request", 13, (byte) 12);
                com.x.dmv2.thriftjava.AcceptMessageRequest.ADAPTER.write(protocol, ((AcceptMessageRequest) struct).m76768getValue());
            } else if (struct instanceof NicknameMessage) {
                protocol.mo14136v3("nickname_message", 14, (byte) 12);
                com.x.dmv2.thriftjava.NicknameMessage.ADAPTER.write(protocol, ((NicknameMessage) struct).m76777getValue());
            } else if (struct instanceof SetVerifiedStatus) {
                protocol.mo14136v3("set_verified_status", 15, (byte) 12);
                com.x.dmv2.thriftjava.SetVerifiedStatus.ADAPTER.write(protocol, ((SetVerifiedStatus) struct).m76782getValue());
            } else if (struct instanceof AvCallStarted) {
                protocol.mo14136v3("av_call_started", 16, (byte) 12);
                AVCallStarted.ADAPTER.write(protocol, ((AvCallStarted) struct).m76771getValue());
            } else if (!(struct instanceof Unknown)) {
                throw new NoWhenBranchMatchedException();
            }
            protocol.mo14134i0();
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$NicknameMessage;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/NicknameMessage;", "<init>", "(Lcom/x/dmv2/thriftjava/NicknameMessage;)V", "getValue", "()Lcom/x/dmv2/thriftjava/NicknameMessage;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class NicknameMessage extends MessageEntryContents {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.NicknameMessage value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public NicknameMessage(@InterfaceC88464a com.x.dmv2.thriftjava.NicknameMessage value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ NicknameMessage copy$default(NicknameMessage nicknameMessage, com.x.dmv2.thriftjava.NicknameMessage nicknameMessage2, int i, Object obj) {
            if ((i & 1) != 0) {
                nicknameMessage2 = nicknameMessage.value;
            }
            return nicknameMessage.copy(nicknameMessage2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.NicknameMessage getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final NicknameMessage copy(@InterfaceC88464a com.x.dmv2.thriftjava.NicknameMessage value) {
            Intrinsics.m65272h(value, "value");
            return new NicknameMessage(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof NicknameMessage) && Intrinsics.m65267c(this.value, ((NicknameMessage) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.NicknameMessage m76777getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(nickname_message=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$PinConversation;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/PinConversation;", "<init>", "(Lcom/x/dmv2/thriftjava/PinConversation;)V", "getValue", "()Lcom/x/dmv2/thriftjava/PinConversation;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class PinConversation extends MessageEntryContents {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.PinConversation value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public PinConversation(@InterfaceC88464a com.x.dmv2.thriftjava.PinConversation value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ PinConversation copy$default(PinConversation pinConversation, com.x.dmv2.thriftjava.PinConversation pinConversation2, int i, Object obj) {
            if ((i & 1) != 0) {
                pinConversation2 = pinConversation.value;
            }
            return pinConversation.copy(pinConversation2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.PinConversation getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final PinConversation copy(@InterfaceC88464a com.x.dmv2.thriftjava.PinConversation value) {
            Intrinsics.m65272h(value, "value");
            return new PinConversation(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof PinConversation) && Intrinsics.m65267c(this.value, ((PinConversation) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.PinConversation m76778getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(pin_conversation=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$ReactionAdd;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/MessageReactionAdd;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageReactionAdd;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageReactionAdd;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class ReactionAdd extends MessageEntryContents {

        @InterfaceC88464a
        private final MessageReactionAdd value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public ReactionAdd(@InterfaceC88464a MessageReactionAdd value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ ReactionAdd copy$default(ReactionAdd reactionAdd, MessageReactionAdd messageReactionAdd, int i, Object obj) {
            if ((i & 1) != 0) {
                messageReactionAdd = reactionAdd.value;
            }
            return reactionAdd.copy(messageReactionAdd);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final MessageReactionAdd getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final ReactionAdd copy(@InterfaceC88464a MessageReactionAdd value) {
            Intrinsics.m65272h(value, "value");
            return new ReactionAdd(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof ReactionAdd) && Intrinsics.m65267c(this.value, ((ReactionAdd) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final MessageReactionAdd m76779getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(reaction_add=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$ReactionRemove;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/MessageReactionRemove;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageReactionRemove;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageReactionRemove;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class ReactionRemove extends MessageEntryContents {

        @InterfaceC88464a
        private final MessageReactionRemove value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public ReactionRemove(@InterfaceC88464a MessageReactionRemove value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ ReactionRemove copy$default(ReactionRemove reactionRemove, MessageReactionRemove messageReactionRemove, int i, Object obj) {
            if ((i & 1) != 0) {
                messageReactionRemove = reactionRemove.value;
            }
            return reactionRemove.copy(messageReactionRemove);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final MessageReactionRemove getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final ReactionRemove copy(@InterfaceC88464a MessageReactionRemove value) {
            Intrinsics.m65272h(value, "value");
            return new ReactionRemove(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof ReactionRemove) && Intrinsics.m65267c(this.value, ((ReactionRemove) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final MessageReactionRemove m76780getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(reaction_remove=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$ScreenCaptureDetected;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/ScreenCaptureDetected;", "<init>", "(Lcom/x/dmv2/thriftjava/ScreenCaptureDetected;)V", "getValue", "()Lcom/x/dmv2/thriftjava/ScreenCaptureDetected;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class ScreenCaptureDetected extends MessageEntryContents {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.ScreenCaptureDetected value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public ScreenCaptureDetected(@InterfaceC88464a com.x.dmv2.thriftjava.ScreenCaptureDetected value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ ScreenCaptureDetected copy$default(ScreenCaptureDetected screenCaptureDetected, com.x.dmv2.thriftjava.ScreenCaptureDetected screenCaptureDetected2, int i, Object obj) {
            if ((i & 1) != 0) {
                screenCaptureDetected2 = screenCaptureDetected.value;
            }
            return screenCaptureDetected.copy(screenCaptureDetected2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.ScreenCaptureDetected getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final ScreenCaptureDetected copy(@InterfaceC88464a com.x.dmv2.thriftjava.ScreenCaptureDetected value) {
            Intrinsics.m65272h(value, "value");
            return new ScreenCaptureDetected(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof ScreenCaptureDetected) && Intrinsics.m65267c(this.value, ((ScreenCaptureDetected) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.ScreenCaptureDetected m76781getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(screen_capture_detected=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$SetVerifiedStatus;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/SetVerifiedStatus;", "<init>", "(Lcom/x/dmv2/thriftjava/SetVerifiedStatus;)V", "getValue", "()Lcom/x/dmv2/thriftjava/SetVerifiedStatus;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class SetVerifiedStatus extends MessageEntryContents {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.SetVerifiedStatus value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public SetVerifiedStatus(@InterfaceC88464a com.x.dmv2.thriftjava.SetVerifiedStatus value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ SetVerifiedStatus copy$default(SetVerifiedStatus setVerifiedStatus, com.x.dmv2.thriftjava.SetVerifiedStatus setVerifiedStatus2, int i, Object obj) {
            if ((i & 1) != 0) {
                setVerifiedStatus2 = setVerifiedStatus.value;
            }
            return setVerifiedStatus.copy(setVerifiedStatus2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.SetVerifiedStatus getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final SetVerifiedStatus copy(@InterfaceC88464a com.x.dmv2.thriftjava.SetVerifiedStatus value) {
            Intrinsics.m65272h(value, "value");
            return new SetVerifiedStatus(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof SetVerifiedStatus) && Intrinsics.m65267c(this.value, ((SetVerifiedStatus) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.SetVerifiedStatus m76782getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(set_verified_status=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\n\u0002\u0010\u000e\n\u0000\bÆ\n\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0013\u0010\u0004\u001a\u00020\u00052\b\u0010\u0006\u001a\u0004\u0018\u00010\u0007HÖ\u0003J\t\u0010\b\u001a\u00020\tHÖ\u0001J\t\u0010\n\u001a\u00020\u000bHÖ\u0001¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$Unknown;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "<init>", "()V", "equals", "", "other", "", "hashCode", "", "toString", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Unknown extends MessageEntryContents {

        @InterfaceC88464a
        public static final Unknown INSTANCE = new Unknown();

        private Unknown() {
            super(null);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            return this == other || (other instanceof Unknown);
        }

        public int hashCode() {
            return -1731007602;
        }

        @InterfaceC88464a
        public String toString() {
            return "Unknown";
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEntryContents$UnpinConversation;", "Lcom/x/dmv2/thriftjava/MessageEntryContents;", "value", "Lcom/x/dmv2/thriftjava/UnpinConversation;", "<init>", "(Lcom/x/dmv2/thriftjava/UnpinConversation;)V", "getValue", "()Lcom/x/dmv2/thriftjava/UnpinConversation;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class UnpinConversation extends MessageEntryContents {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.UnpinConversation value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public UnpinConversation(@InterfaceC88464a com.x.dmv2.thriftjava.UnpinConversation value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ UnpinConversation copy$default(UnpinConversation unpinConversation, com.x.dmv2.thriftjava.UnpinConversation unpinConversation2, int i, Object obj) {
            if ((i & 1) != 0) {
                unpinConversation2 = unpinConversation.value;
            }
            return unpinConversation.copy(unpinConversation2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.UnpinConversation getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final UnpinConversation copy(@InterfaceC88464a com.x.dmv2.thriftjava.UnpinConversation value) {
            Intrinsics.m65272h(value, "value");
            return new UnpinConversation(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof UnpinConversation) && Intrinsics.m65267c(this.value, ((UnpinConversation) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.UnpinConversation m76783getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageEntryContents(unpin_conversation=" + this.value + Separators.RPAREN;
        }
    }

    public /* synthetic */ MessageEntryContents(DefaultConstructorMarker defaultConstructorMarker) {
        this();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }

    private MessageEntryContents() {
    }
}