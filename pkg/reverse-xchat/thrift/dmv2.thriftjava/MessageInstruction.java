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

@Metadata(m64929d1 = {"\u00008\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u000b\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\b6\u0018\u0000 \t2\u00020\u0001:\t\n\u000b\f\r\u000e\u000f\u0010\u0011\tB\t\b\u0004¢\u0006\u0004\b\u0002\u0010\u0003J\u0017\u0010\u0007\u001a\u00020\u00062\u0006\u0010\u0005\u001a\u00020\u0004H\u0016¢\u0006\u0004\b\u0007\u0010\b\u0082\u0001\u0007\u0012\u0013\u0014\u0015\u0016\u0017\u0018¨\u0006\u0019"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageInstruction;", "Lcom/bendb/thrifty/a;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "Companion", "PullMessagesInstruction", "KeepAliveInstruction", "PullMessagesFinishedInstruction", "PinReminderInstruction", "SwitchToHybridPullInstruction", "DisplayTemporaryPasscodeInstruction", "Unknown", "MessageInstructionAdapter", "Lcom/x/dmv2/thriftjava/MessageInstruction$DisplayTemporaryPasscodeInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction$KeepAliveInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction$PinReminderInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction$PullMessagesFinishedInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction$PullMessagesInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction$SwitchToHybridPullInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction$Unknown;", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public abstract class MessageInstruction implements InterfaceC11261a {

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MessageInstructionAdapter();

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageInstruction$DisplayTemporaryPasscodeInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction;", "value", "Lcom/x/dmv2/thriftjava/DisplayTemporaryPasscodeInstruction;", "<init>", "(Lcom/x/dmv2/thriftjava/DisplayTemporaryPasscodeInstruction;)V", "getValue", "()Lcom/x/dmv2/thriftjava/DisplayTemporaryPasscodeInstruction;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class DisplayTemporaryPasscodeInstruction extends MessageInstruction {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.DisplayTemporaryPasscodeInstruction value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public DisplayTemporaryPasscodeInstruction(@InterfaceC88464a com.x.dmv2.thriftjava.DisplayTemporaryPasscodeInstruction value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ DisplayTemporaryPasscodeInstruction copy$default(DisplayTemporaryPasscodeInstruction displayTemporaryPasscodeInstruction, com.x.dmv2.thriftjava.DisplayTemporaryPasscodeInstruction displayTemporaryPasscodeInstruction2, int i, Object obj) {
            if ((i & 1) != 0) {
                displayTemporaryPasscodeInstruction2 = displayTemporaryPasscodeInstruction.value;
            }
            return displayTemporaryPasscodeInstruction.copy(displayTemporaryPasscodeInstruction2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.DisplayTemporaryPasscodeInstruction getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final DisplayTemporaryPasscodeInstruction copy(@InterfaceC88464a com.x.dmv2.thriftjava.DisplayTemporaryPasscodeInstruction value) {
            Intrinsics.m65272h(value, "value");
            return new DisplayTemporaryPasscodeInstruction(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof DisplayTemporaryPasscodeInstruction) && Intrinsics.m65267c(this.value, ((DisplayTemporaryPasscodeInstruction) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.DisplayTemporaryPasscodeInstruction m76797getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageInstruction(displayTemporaryPasscodeInstruction=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageInstruction$KeepAliveInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction;", "value", "Lcom/x/dmv2/thriftjava/KeepAliveInstruction;", "<init>", "(Lcom/x/dmv2/thriftjava/KeepAliveInstruction;)V", "getValue", "()Lcom/x/dmv2/thriftjava/KeepAliveInstruction;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class KeepAliveInstruction extends MessageInstruction {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.KeepAliveInstruction value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public KeepAliveInstruction(@InterfaceC88464a com.x.dmv2.thriftjava.KeepAliveInstruction value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ KeepAliveInstruction copy$default(KeepAliveInstruction keepAliveInstruction, com.x.dmv2.thriftjava.KeepAliveInstruction keepAliveInstruction2, int i, Object obj) {
            if ((i & 1) != 0) {
                keepAliveInstruction2 = keepAliveInstruction.value;
            }
            return keepAliveInstruction.copy(keepAliveInstruction2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.KeepAliveInstruction getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final KeepAliveInstruction copy(@InterfaceC88464a com.x.dmv2.thriftjava.KeepAliveInstruction value) {
            Intrinsics.m65272h(value, "value");
            return new KeepAliveInstruction(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof KeepAliveInstruction) && Intrinsics.m65267c(this.value, ((KeepAliveInstruction) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.KeepAliveInstruction m76798getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageInstruction(keepAliveInstruction=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageInstruction$MessageInstructionAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MessageInstruction;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MessageInstruction;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MessageInstruction;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MessageInstructionAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MessageInstruction m85661read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            MessageInstruction pullMessagesInstruction;
            Intrinsics.m65272h(protocol, "protocol");
            MessageInstruction messageInstruction = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    if (messageInstruction != null) {
                        return messageInstruction;
                    }
                    throw new IllegalStateException("unreadable");
                }
                switch (c11265cMo14127V2.f38393b) {
                    case 1:
                        if (b == 12) {
                            pullMessagesInstruction = new PullMessagesInstruction((com.x.dmv2.thriftjava.PullMessagesInstruction) com.x.dmv2.thriftjava.PullMessagesInstruction.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 2:
                        if (b == 12) {
                            pullMessagesInstruction = new KeepAliveInstruction((com.x.dmv2.thriftjava.KeepAliveInstruction) com.x.dmv2.thriftjava.KeepAliveInstruction.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 3:
                        if (b == 12) {
                            pullMessagesInstruction = new PullMessagesFinishedInstruction((com.x.dmv2.thriftjava.PullMessagesFinishedInstruction) com.x.dmv2.thriftjava.PullMessagesFinishedInstruction.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 4:
                        if (b == 12) {
                            pullMessagesInstruction = new PinReminderInstruction((com.x.dmv2.thriftjava.PinReminderInstruction) com.x.dmv2.thriftjava.PinReminderInstruction.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 5:
                        if (b == 12) {
                            pullMessagesInstruction = new SwitchToHybridPullInstruction((com.x.dmv2.thriftjava.SwitchToHybridPullInstruction) com.x.dmv2.thriftjava.SwitchToHybridPullInstruction.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    case 6:
                        if (b == 12) {
                            pullMessagesInstruction = new DisplayTemporaryPasscodeInstruction((com.x.dmv2.thriftjava.DisplayTemporaryPasscodeInstruction) com.x.dmv2.thriftjava.DisplayTemporaryPasscodeInstruction.ADAPTER.read(protocol));
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    default:
                        messageInstruction = Unknown.INSTANCE;
                        C11272a.m14141a(protocol, b);
                        continue;
                }
                messageInstruction = pullMessagesInstruction;
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MessageInstruction struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MessageInstruction");
            if (struct instanceof PullMessagesInstruction) {
                protocol.mo14136v3("pullMessagesInstruction", 1, (byte) 12);
                com.x.dmv2.thriftjava.PullMessagesInstruction.ADAPTER.write(protocol, ((PullMessagesInstruction) struct).m76801getValue());
            } else if (struct instanceof KeepAliveInstruction) {
                protocol.mo14136v3("keepAliveInstruction", 2, (byte) 12);
                com.x.dmv2.thriftjava.KeepAliveInstruction.ADAPTER.write(protocol, ((KeepAliveInstruction) struct).m76798getValue());
            } else if (struct instanceof PullMessagesFinishedInstruction) {
                protocol.mo14136v3("pullMessagesFinishedInstruction", 3, (byte) 12);
                com.x.dmv2.thriftjava.PullMessagesFinishedInstruction.ADAPTER.write(protocol, ((PullMessagesFinishedInstruction) struct).m76800getValue());
            } else if (struct instanceof PinReminderInstruction) {
                protocol.mo14136v3("pinReminderInstruction", 4, (byte) 12);
                com.x.dmv2.thriftjava.PinReminderInstruction.ADAPTER.write(protocol, ((PinReminderInstruction) struct).m76799getValue());
            } else if (struct instanceof SwitchToHybridPullInstruction) {
                protocol.mo14136v3("switchToHybridPullInstruction", 5, (byte) 12);
                com.x.dmv2.thriftjava.SwitchToHybridPullInstruction.ADAPTER.write(protocol, ((SwitchToHybridPullInstruction) struct).m76802getValue());
            } else if (struct instanceof DisplayTemporaryPasscodeInstruction) {
                protocol.mo14136v3("displayTemporaryPasscodeInstruction", 6, (byte) 12);
                com.x.dmv2.thriftjava.DisplayTemporaryPasscodeInstruction.ADAPTER.write(protocol, ((DisplayTemporaryPasscodeInstruction) struct).m76797getValue());
            } else if (!(struct instanceof Unknown)) {
                throw new NoWhenBranchMatchedException();
            }
            protocol.mo14134i0();
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageInstruction$PinReminderInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction;", "value", "Lcom/x/dmv2/thriftjava/PinReminderInstruction;", "<init>", "(Lcom/x/dmv2/thriftjava/PinReminderInstruction;)V", "getValue", "()Lcom/x/dmv2/thriftjava/PinReminderInstruction;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class PinReminderInstruction extends MessageInstruction {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.PinReminderInstruction value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public PinReminderInstruction(@InterfaceC88464a com.x.dmv2.thriftjava.PinReminderInstruction value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ PinReminderInstruction copy$default(PinReminderInstruction pinReminderInstruction, com.x.dmv2.thriftjava.PinReminderInstruction pinReminderInstruction2, int i, Object obj) {
            if ((i & 1) != 0) {
                pinReminderInstruction2 = pinReminderInstruction.value;
            }
            return pinReminderInstruction.copy(pinReminderInstruction2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.PinReminderInstruction getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final PinReminderInstruction copy(@InterfaceC88464a com.x.dmv2.thriftjava.PinReminderInstruction value) {
            Intrinsics.m65272h(value, "value");
            return new PinReminderInstruction(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof PinReminderInstruction) && Intrinsics.m65267c(this.value, ((PinReminderInstruction) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.PinReminderInstruction m76799getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageInstruction(pinReminderInstruction=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageInstruction$PullMessagesFinishedInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction;", "value", "Lcom/x/dmv2/thriftjava/PullMessagesFinishedInstruction;", "<init>", "(Lcom/x/dmv2/thriftjava/PullMessagesFinishedInstruction;)V", "getValue", "()Lcom/x/dmv2/thriftjava/PullMessagesFinishedInstruction;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class PullMessagesFinishedInstruction extends MessageInstruction {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.PullMessagesFinishedInstruction value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public PullMessagesFinishedInstruction(@InterfaceC88464a com.x.dmv2.thriftjava.PullMessagesFinishedInstruction value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ PullMessagesFinishedInstruction copy$default(PullMessagesFinishedInstruction pullMessagesFinishedInstruction, com.x.dmv2.thriftjava.PullMessagesFinishedInstruction pullMessagesFinishedInstruction2, int i, Object obj) {
            if ((i & 1) != 0) {
                pullMessagesFinishedInstruction2 = pullMessagesFinishedInstruction.value;
            }
            return pullMessagesFinishedInstruction.copy(pullMessagesFinishedInstruction2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.PullMessagesFinishedInstruction getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final PullMessagesFinishedInstruction copy(@InterfaceC88464a com.x.dmv2.thriftjava.PullMessagesFinishedInstruction value) {
            Intrinsics.m65272h(value, "value");
            return new PullMessagesFinishedInstruction(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof PullMessagesFinishedInstruction) && Intrinsics.m65267c(this.value, ((PullMessagesFinishedInstruction) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.PullMessagesFinishedInstruction m76800getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageInstruction(pullMessagesFinishedInstruction=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageInstruction$PullMessagesInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction;", "value", "Lcom/x/dmv2/thriftjava/PullMessagesInstruction;", "<init>", "(Lcom/x/dmv2/thriftjava/PullMessagesInstruction;)V", "getValue", "()Lcom/x/dmv2/thriftjava/PullMessagesInstruction;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class PullMessagesInstruction extends MessageInstruction {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.PullMessagesInstruction value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public PullMessagesInstruction(@InterfaceC88464a com.x.dmv2.thriftjava.PullMessagesInstruction value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ PullMessagesInstruction copy$default(PullMessagesInstruction pullMessagesInstruction, com.x.dmv2.thriftjava.PullMessagesInstruction pullMessagesInstruction2, int i, Object obj) {
            if ((i & 1) != 0) {
                pullMessagesInstruction2 = pullMessagesInstruction.value;
            }
            return pullMessagesInstruction.copy(pullMessagesInstruction2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.PullMessagesInstruction getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final PullMessagesInstruction copy(@InterfaceC88464a com.x.dmv2.thriftjava.PullMessagesInstruction value) {
            Intrinsics.m65272h(value, "value");
            return new PullMessagesInstruction(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof PullMessagesInstruction) && Intrinsics.m65267c(this.value, ((PullMessagesInstruction) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.PullMessagesInstruction m76801getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageInstruction(pullMessagesInstruction=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageInstruction$SwitchToHybridPullInstruction;", "Lcom/x/dmv2/thriftjava/MessageInstruction;", "value", "Lcom/x/dmv2/thriftjava/SwitchToHybridPullInstruction;", "<init>", "(Lcom/x/dmv2/thriftjava/SwitchToHybridPullInstruction;)V", "getValue", "()Lcom/x/dmv2/thriftjava/SwitchToHybridPullInstruction;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class SwitchToHybridPullInstruction extends MessageInstruction {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.SwitchToHybridPullInstruction value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public SwitchToHybridPullInstruction(@InterfaceC88464a com.x.dmv2.thriftjava.SwitchToHybridPullInstruction value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ SwitchToHybridPullInstruction copy$default(SwitchToHybridPullInstruction switchToHybridPullInstruction, com.x.dmv2.thriftjava.SwitchToHybridPullInstruction switchToHybridPullInstruction2, int i, Object obj) {
            if ((i & 1) != 0) {
                switchToHybridPullInstruction2 = switchToHybridPullInstruction.value;
            }
            return switchToHybridPullInstruction.copy(switchToHybridPullInstruction2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.SwitchToHybridPullInstruction getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final SwitchToHybridPullInstruction copy(@InterfaceC88464a com.x.dmv2.thriftjava.SwitchToHybridPullInstruction value) {
            Intrinsics.m65272h(value, "value");
            return new SwitchToHybridPullInstruction(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof SwitchToHybridPullInstruction) && Intrinsics.m65267c(this.value, ((SwitchToHybridPullInstruction) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.SwitchToHybridPullInstruction m76802getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "MessageInstruction(switchToHybridPullInstruction=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\n\u0002\u0010\u000e\n\u0000\bÆ\n\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0013\u0010\u0004\u001a\u00020\u00052\b\u0010\u0006\u001a\u0004\u0018\u00010\u0007HÖ\u0003J\t\u0010\b\u001a\u00020\tHÖ\u0001J\t\u0010\n\u001a\u00020\u000bHÖ\u0001¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageInstruction$Unknown;", "Lcom/x/dmv2/thriftjava/MessageInstruction;", "<init>", "()V", "equals", "", "other", "", "hashCode", "", "toString", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Unknown extends MessageInstruction {

        @InterfaceC88464a
        public static final Unknown INSTANCE = new Unknown();

        private Unknown() {
            super(null);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            return this == other || (other instanceof Unknown);
        }

        public int hashCode() {
            return -2132487696;
        }

        @InterfaceC88464a
        public String toString() {
            return "Unknown";
        }
    }

    public /* synthetic */ MessageInstruction(DefaultConstructorMarker defaultConstructorMarker) {
        this();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }

    private MessageInstruction() {
    }
}