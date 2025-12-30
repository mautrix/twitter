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

@Metadata(m64929d1 = {"\u0000@\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\r\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\b6\u0018\u0000 \t2\u00020\u0001:\u000b\n\u000b\f\r\u000e\u000f\u0010\u0011\u0012\u0013\tB\t\b\u0004¢\u0006\u0004\b\u0002\u0010\u0003J\u0017\u0010\u0007\u001a\u00020\u00062\u0006\u0010\u0005\u001a\u00020\u0004H\u0016¢\u0006\u0004\b\u0007\u0010\b\u0082\u0001\t\u0014\u0015\u0016\u0017\u0018\u0019\u001a\u001b\u001c¨\u0006\u001d"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "Lcom/bendb/thrifty/a;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "Companion", "MessageDurationChange", "MessageDurationRemove", "MuteConversation", "UnmuteConversation", "EnableScreenCaptureDetection", "DisableScreenCaptureDetection", "EnableScreenCaptureBlocking", "DisableScreenCaptureBlocking", "Unknown", "ConversationMetadataChangeAdapter", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange$DisableScreenCaptureBlocking;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange$DisableScreenCaptureDetection;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange$EnableScreenCaptureBlocking;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange$EnableScreenCaptureDetection;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange$MessageDurationChange;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange$MessageDurationRemove;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange$MuteConversation;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange$Unknown;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange$UnmuteConversation;", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public abstract class ConversationMetadataChange implements InterfaceC11261a {

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new ConversationMetadataChangeAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationMetadataChange$ConversationMetadataChangeAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/ConversationMetadataChange;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class ConversationMetadataChangeAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public ConversationMetadataChange m83700read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            ConversationMetadataChange messageDurationChange;
            Intrinsics.m65272h(protocol, "protocol");
            ConversationMetadataChange conversationMetadataChange = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    if (conversationMetadataChange != null) {
                        return conversationMetadataChange;
                    }
                    throw new IllegalStateException("unreadable");
                }
                switch (c11265cMo14127V2.f38393b) {
                    case 1:
                        if (b == 12) {
                            messageDurationChange = new MessageDurationChange((com.x.dmv2.thriftjava.MessageDurationChange) com.x.dmv2.thriftjava.MessageDurationChange.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 2:
                        if (b == 12) {
                            messageDurationChange = new MessageDurationRemove((com.x.dmv2.thriftjava.MessageDurationRemove) com.x.dmv2.thriftjava.MessageDurationRemove.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 3:
                        if (b == 12) {
                            messageDurationChange = new MuteConversation((com.x.dmv2.thriftjava.MuteConversation) com.x.dmv2.thriftjava.MuteConversation.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 4:
                        if (b == 12) {
                            messageDurationChange = new UnmuteConversation((com.x.dmv2.thriftjava.UnmuteConversation) com.x.dmv2.thriftjava.UnmuteConversation.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 5:
                        if (b == 12) {
                            messageDurationChange = new EnableScreenCaptureDetection((com.x.dmv2.thriftjava.EnableScreenCaptureDetection) com.x.dmv2.thriftjava.EnableScreenCaptureDetection.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 6:
                        if (b == 12) {
                            messageDurationChange = new DisableScreenCaptureDetection((com.x.dmv2.thriftjava.DisableScreenCaptureDetection) com.x.dmv2.thriftjava.DisableScreenCaptureDetection.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 7:
                        if (b == 12) {
                            messageDurationChange = new EnableScreenCaptureBlocking((com.x.dmv2.thriftjava.EnableScreenCaptureBlocking) com.x.dmv2.thriftjava.EnableScreenCaptureBlocking.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 8:
                        if (b == 12) {
                            messageDurationChange = new DisableScreenCaptureBlocking((com.x.dmv2.thriftjava.DisableScreenCaptureBlocking) com.x.dmv2.thriftjava.DisableScreenCaptureBlocking.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    default:
                        conversationMetadataChange = Unknown.INSTANCE;
                        C11272a.m14141a(protocol, b);
                        continue;
                }
                conversationMetadataChange = messageDurationChange;
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a ConversationMetadataChange struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("ConversationMetadataChange");
            if (struct instanceof MessageDurationChange) {
                protocol.mo14136v3("message_duration_change", 1, (byte) 12);
                com.x.dmv2.thriftjava.MessageDurationChange.ADAPTER.write(protocol, ((MessageDurationChange) struct).m76743getValue());
            } else if (struct instanceof MessageDurationRemove) {
                protocol.mo14136v3("message_duration_remove", 2, (byte) 12);
                com.x.dmv2.thriftjava.MessageDurationRemove.ADAPTER.write(protocol, ((MessageDurationRemove) struct).m76744getValue());
            } else if (struct instanceof MuteConversation) {
                protocol.mo14136v3("mute_conversation", 3, (byte) 12);
                com.x.dmv2.thriftjava.MuteConversation.ADAPTER.write(protocol, ((MuteConversation) struct).m76745getValue());
            } else if (struct instanceof UnmuteConversation) {
                protocol.mo14136v3("unmute_conversation", 4, (byte) 12);
                com.x.dmv2.thriftjava.UnmuteConversation.ADAPTER.write(protocol, ((UnmuteConversation) struct).m76746getValue());
            } else if (struct instanceof EnableScreenCaptureDetection) {
                protocol.mo14136v3("enable_screen_capture_detection", 5, (byte) 12);
                com.x.dmv2.thriftjava.EnableScreenCaptureDetection.ADAPTER.write(protocol, ((EnableScreenCaptureDetection) struct).m76742getValue());
            } else if (struct instanceof DisableScreenCaptureDetection) {
                protocol.mo14136v3("disable_screen_capture_detection", 6, (byte) 12);
                com.x.dmv2.thriftjava.DisableScreenCaptureDetection.ADAPTER.write(protocol, ((DisableScreenCaptureDetection) struct).m76740getValue());
            } else if (struct instanceof EnableScreenCaptureBlocking) {
                protocol.mo14136v3("enable_screen_capture_blocking", 7, (byte) 12);
                com.x.dmv2.thriftjava.EnableScreenCaptureBlocking.ADAPTER.write(protocol, ((EnableScreenCaptureBlocking) struct).m76741getValue());
            } else if (struct instanceof DisableScreenCaptureBlocking) {
                protocol.mo14136v3("disable_screen_capture_blocking", 8, (byte) 12);
                com.x.dmv2.thriftjava.DisableScreenCaptureBlocking.ADAPTER.write(protocol, ((DisableScreenCaptureBlocking) struct).m76739getValue());
            } else if (!(struct instanceof Unknown)) {
                throw new NoWhenBranchMatchedException();
            }
            protocol.mo14134i0();
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationMetadataChange$DisableScreenCaptureBlocking;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "value", "Lcom/x/dmv2/thriftjava/DisableScreenCaptureBlocking;", "<init>", "(Lcom/x/dmv2/thriftjava/DisableScreenCaptureBlocking;)V", "getValue", "()Lcom/x/dmv2/thriftjava/DisableScreenCaptureBlocking;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class DisableScreenCaptureBlocking extends ConversationMetadataChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.DisableScreenCaptureBlocking value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public DisableScreenCaptureBlocking(@InterfaceC88464a com.x.dmv2.thriftjava.DisableScreenCaptureBlocking value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ DisableScreenCaptureBlocking copy$default(DisableScreenCaptureBlocking disableScreenCaptureBlocking, com.x.dmv2.thriftjava.DisableScreenCaptureBlocking disableScreenCaptureBlocking2, int i, Object obj) {
            if ((i & 1) != 0) {
                disableScreenCaptureBlocking2 = disableScreenCaptureBlocking.value;
            }
            return disableScreenCaptureBlocking.copy(disableScreenCaptureBlocking2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.DisableScreenCaptureBlocking getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final DisableScreenCaptureBlocking copy(@InterfaceC88464a com.x.dmv2.thriftjava.DisableScreenCaptureBlocking value) {
            Intrinsics.m65272h(value, "value");
            return new DisableScreenCaptureBlocking(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof DisableScreenCaptureBlocking) && Intrinsics.m65267c(this.value, ((DisableScreenCaptureBlocking) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.DisableScreenCaptureBlocking m76739getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "ConversationMetadataChange(disable_screen_capture_blocking=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationMetadataChange$DisableScreenCaptureDetection;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "value", "Lcom/x/dmv2/thriftjava/DisableScreenCaptureDetection;", "<init>", "(Lcom/x/dmv2/thriftjava/DisableScreenCaptureDetection;)V", "getValue", "()Lcom/x/dmv2/thriftjava/DisableScreenCaptureDetection;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class DisableScreenCaptureDetection extends ConversationMetadataChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.DisableScreenCaptureDetection value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public DisableScreenCaptureDetection(@InterfaceC88464a com.x.dmv2.thriftjava.DisableScreenCaptureDetection value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ DisableScreenCaptureDetection copy$default(DisableScreenCaptureDetection disableScreenCaptureDetection, com.x.dmv2.thriftjava.DisableScreenCaptureDetection disableScreenCaptureDetection2, int i, Object obj) {
            if ((i & 1) != 0) {
                disableScreenCaptureDetection2 = disableScreenCaptureDetection.value;
            }
            return disableScreenCaptureDetection.copy(disableScreenCaptureDetection2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.DisableScreenCaptureDetection getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final DisableScreenCaptureDetection copy(@InterfaceC88464a com.x.dmv2.thriftjava.DisableScreenCaptureDetection value) {
            Intrinsics.m65272h(value, "value");
            return new DisableScreenCaptureDetection(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof DisableScreenCaptureDetection) && Intrinsics.m65267c(this.value, ((DisableScreenCaptureDetection) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.DisableScreenCaptureDetection m76740getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "ConversationMetadataChange(disable_screen_capture_detection=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationMetadataChange$EnableScreenCaptureBlocking;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "value", "Lcom/x/dmv2/thriftjava/EnableScreenCaptureBlocking;", "<init>", "(Lcom/x/dmv2/thriftjava/EnableScreenCaptureBlocking;)V", "getValue", "()Lcom/x/dmv2/thriftjava/EnableScreenCaptureBlocking;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class EnableScreenCaptureBlocking extends ConversationMetadataChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.EnableScreenCaptureBlocking value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public EnableScreenCaptureBlocking(@InterfaceC88464a com.x.dmv2.thriftjava.EnableScreenCaptureBlocking value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ EnableScreenCaptureBlocking copy$default(EnableScreenCaptureBlocking enableScreenCaptureBlocking, com.x.dmv2.thriftjava.EnableScreenCaptureBlocking enableScreenCaptureBlocking2, int i, Object obj) {
            if ((i & 1) != 0) {
                enableScreenCaptureBlocking2 = enableScreenCaptureBlocking.value;
            }
            return enableScreenCaptureBlocking.copy(enableScreenCaptureBlocking2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.EnableScreenCaptureBlocking getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final EnableScreenCaptureBlocking copy(@InterfaceC88464a com.x.dmv2.thriftjava.EnableScreenCaptureBlocking value) {
            Intrinsics.m65272h(value, "value");
            return new EnableScreenCaptureBlocking(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof EnableScreenCaptureBlocking) && Intrinsics.m65267c(this.value, ((EnableScreenCaptureBlocking) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.EnableScreenCaptureBlocking m76741getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "ConversationMetadataChange(enable_screen_capture_blocking=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationMetadataChange$EnableScreenCaptureDetection;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "value", "Lcom/x/dmv2/thriftjava/EnableScreenCaptureDetection;", "<init>", "(Lcom/x/dmv2/thriftjava/EnableScreenCaptureDetection;)V", "getValue", "()Lcom/x/dmv2/thriftjava/EnableScreenCaptureDetection;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class EnableScreenCaptureDetection extends ConversationMetadataChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.EnableScreenCaptureDetection value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public EnableScreenCaptureDetection(@InterfaceC88464a com.x.dmv2.thriftjava.EnableScreenCaptureDetection value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ EnableScreenCaptureDetection copy$default(EnableScreenCaptureDetection enableScreenCaptureDetection, com.x.dmv2.thriftjava.EnableScreenCaptureDetection enableScreenCaptureDetection2, int i, Object obj) {
            if ((i & 1) != 0) {
                enableScreenCaptureDetection2 = enableScreenCaptureDetection.value;
            }
            return enableScreenCaptureDetection.copy(enableScreenCaptureDetection2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.EnableScreenCaptureDetection getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final EnableScreenCaptureDetection copy(@InterfaceC88464a com.x.dmv2.thriftjava.EnableScreenCaptureDetection value) {
            Intrinsics.m65272h(value, "value");
            return new EnableScreenCaptureDetection(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof EnableScreenCaptureDetection) && Intrinsics.m65267c(this.value, ((EnableScreenCaptureDetection) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.EnableScreenCaptureDetection m76742getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "ConversationMetadataChange(enable_screen_capture_detection=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationMetadataChange$MessageDurationChange;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "value", "Lcom/x/dmv2/thriftjava/MessageDurationChange;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageDurationChange;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageDurationChange;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MessageDurationChange extends ConversationMetadataChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MessageDurationChange value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MessageDurationChange(@InterfaceC88464a com.x.dmv2.thriftjava.MessageDurationChange value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MessageDurationChange copy$default(MessageDurationChange messageDurationChange, com.x.dmv2.thriftjava.MessageDurationChange messageDurationChange2, int i, Object obj) {
            if ((i & 1) != 0) {
                messageDurationChange2 = messageDurationChange.value;
            }
            return messageDurationChange.copy(messageDurationChange2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MessageDurationChange getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MessageDurationChange copy(@InterfaceC88464a com.x.dmv2.thriftjava.MessageDurationChange value) {
            Intrinsics.m65272h(value, "value");
            return new MessageDurationChange(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MessageDurationChange) && Intrinsics.m65267c(this.value, ((MessageDurationChange) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MessageDurationChange m76743getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "ConversationMetadataChange(message_duration_change=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationMetadataChange$MessageDurationRemove;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "value", "Lcom/x/dmv2/thriftjava/MessageDurationRemove;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageDurationRemove;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageDurationRemove;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MessageDurationRemove extends ConversationMetadataChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MessageDurationRemove value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MessageDurationRemove(@InterfaceC88464a com.x.dmv2.thriftjava.MessageDurationRemove value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MessageDurationRemove copy$default(MessageDurationRemove messageDurationRemove, com.x.dmv2.thriftjava.MessageDurationRemove messageDurationRemove2, int i, Object obj) {
            if ((i & 1) != 0) {
                messageDurationRemove2 = messageDurationRemove.value;
            }
            return messageDurationRemove.copy(messageDurationRemove2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MessageDurationRemove getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MessageDurationRemove copy(@InterfaceC88464a com.x.dmv2.thriftjava.MessageDurationRemove value) {
            Intrinsics.m65272h(value, "value");
            return new MessageDurationRemove(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MessageDurationRemove) && Intrinsics.m65267c(this.value, ((MessageDurationRemove) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MessageDurationRemove m76744getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "ConversationMetadataChange(message_duration_remove=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationMetadataChange$MuteConversation;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "value", "Lcom/x/dmv2/thriftjava/MuteConversation;", "<init>", "(Lcom/x/dmv2/thriftjava/MuteConversation;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MuteConversation;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MuteConversation extends ConversationMetadataChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MuteConversation value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MuteConversation(@InterfaceC88464a com.x.dmv2.thriftjava.MuteConversation value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MuteConversation copy$default(MuteConversation muteConversation, com.x.dmv2.thriftjava.MuteConversation muteConversation2, int i, Object obj) {
            if ((i & 1) != 0) {
                muteConversation2 = muteConversation.value;
            }
            return muteConversation.copy(muteConversation2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MuteConversation getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MuteConversation copy(@InterfaceC88464a com.x.dmv2.thriftjava.MuteConversation value) {
            Intrinsics.m65272h(value, "value");
            return new MuteConversation(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MuteConversation) && Intrinsics.m65267c(this.value, ((MuteConversation) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MuteConversation m76745getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "ConversationMetadataChange(mute_conversation=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\n\u0002\u0010\u000e\n\u0000\bÆ\n\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0013\u0010\u0004\u001a\u00020\u00052\b\u0010\u0006\u001a\u0004\u0018\u00010\u0007HÖ\u0003J\t\u0010\b\u001a\u00020\tHÖ\u0001J\t\u0010\n\u001a\u00020\u000bHÖ\u0001¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationMetadataChange$Unknown;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "<init>", "()V", "equals", "", "other", "", "hashCode", "", "toString", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Unknown extends ConversationMetadataChange {

        @InterfaceC88464a
        public static final Unknown INSTANCE = new Unknown();

        private Unknown() {
            super(null);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            return this == other || (other instanceof Unknown);
        }

        public int hashCode() {
            return 1182176011;
        }

        @InterfaceC88464a
        public String toString() {
            return "Unknown";
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationMetadataChange$UnmuteConversation;", "Lcom/x/dmv2/thriftjava/ConversationMetadataChange;", "value", "Lcom/x/dmv2/thriftjava/UnmuteConversation;", "<init>", "(Lcom/x/dmv2/thriftjava/UnmuteConversation;)V", "getValue", "()Lcom/x/dmv2/thriftjava/UnmuteConversation;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class UnmuteConversation extends ConversationMetadataChange {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.UnmuteConversation value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public UnmuteConversation(@InterfaceC88464a com.x.dmv2.thriftjava.UnmuteConversation value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ UnmuteConversation copy$default(UnmuteConversation unmuteConversation, com.x.dmv2.thriftjava.UnmuteConversation unmuteConversation2, int i, Object obj) {
            if ((i & 1) != 0) {
                unmuteConversation2 = unmuteConversation.value;
            }
            return unmuteConversation.copy(unmuteConversation2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.UnmuteConversation getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final UnmuteConversation copy(@InterfaceC88464a com.x.dmv2.thriftjava.UnmuteConversation value) {
            Intrinsics.m65272h(value, "value");
            return new UnmuteConversation(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof UnmuteConversation) && Intrinsics.m65267c(this.value, ((UnmuteConversation) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.UnmuteConversation m76746getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "ConversationMetadataChange(unmute_conversation=" + this.value + Separators.RPAREN;
        }
    }

    public /* synthetic */ ConversationMetadataChange(DefaultConstructorMarker defaultConstructorMarker) {
        this();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }

    private ConversationMetadataChange() {
    }
}
