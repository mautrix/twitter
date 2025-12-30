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

@Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\b\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\b6\u0018\u0000 \t2\u00020\u0001:\u0006\n\u000b\f\r\u000e\tB\t\b\u0004¢\u0006\u0004\b\u0002\u0010\u0003J\u0017\u0010\u0007\u001a\u00020\u00062\u0006\u0010\u0005\u001a\u00020\u0004H\u0016¢\u0006\u0004\b\u0007\u0010\b\u0082\u0001\u0004\u000f\u0010\u0011\u0012¨\u0006\u0013"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/Message;", "Lcom/bendb/thrifty/a;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "Companion", "MessageEvent", "MessageInstruction", "BatchedMessageEvents", "Unknown", "MessageAdapter", "Lcom/x/dmv2/thriftjava/Message$BatchedMessageEvents;", "Lcom/x/dmv2/thriftjava/Message$MessageEvent;", "Lcom/x/dmv2/thriftjava/Message$MessageInstruction;", "Lcom/x/dmv2/thriftjava/Message$Unknown;", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public abstract class Message implements InterfaceC11261a {

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MessageAdapter();

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/Message$BatchedMessageEvents;", "Lcom/x/dmv2/thriftjava/Message;", "value", "Lcom/x/dmv2/thriftjava/BatchedMessageEvents;", "<init>", "(Lcom/x/dmv2/thriftjava/BatchedMessageEvents;)V", "getValue", "()Lcom/x/dmv2/thriftjava/BatchedMessageEvents;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class BatchedMessageEvents extends Message {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.BatchedMessageEvents value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public BatchedMessageEvents(@InterfaceC88464a com.x.dmv2.thriftjava.BatchedMessageEvents value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ BatchedMessageEvents copy$default(BatchedMessageEvents batchedMessageEvents, com.x.dmv2.thriftjava.BatchedMessageEvents batchedMessageEvents2, int i, Object obj) {
            if ((i & 1) != 0) {
                batchedMessageEvents2 = batchedMessageEvents.value;
            }
            return batchedMessageEvents.copy(batchedMessageEvents2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.BatchedMessageEvents getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final BatchedMessageEvents copy(@InterfaceC88464a com.x.dmv2.thriftjava.BatchedMessageEvents value) {
            Intrinsics.m65272h(value, "value");
            return new BatchedMessageEvents(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof BatchedMessageEvents) && Intrinsics.m65267c(this.value, ((BatchedMessageEvents) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.BatchedMessageEvents m76760getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "Message(batchedMessageEvents=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/Message$MessageAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/Message;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/Message;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/Message;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MessageAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public Message m85652read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Message batchedMessageEvents;
            Intrinsics.m65272h(protocol, "protocol");
            Message message = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    break;
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            message = Unknown.INSTANCE;
                            C11272a.m14141a(protocol, b);
                        } else if (b == 12) {
                            batchedMessageEvents = new BatchedMessageEvents((com.x.dmv2.thriftjava.BatchedMessageEvents) com.x.dmv2.thriftjava.BatchedMessageEvents.ADAPTER.read(protocol));
                            message = batchedMessageEvents;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 12) {
                        batchedMessageEvents = new MessageInstruction((com.x.dmv2.thriftjava.MessageInstruction) com.x.dmv2.thriftjava.MessageInstruction.ADAPTER.read(protocol));
                        message = batchedMessageEvents;
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 12) {
                    batchedMessageEvents = new MessageEvent((com.x.dmv2.thriftjava.MessageEvent) com.x.dmv2.thriftjava.MessageEvent.ADAPTER.read(protocol));
                    message = batchedMessageEvents;
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
            if (message != null) {
                return message;
            }
            throw new IllegalStateException("unreadable");
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a Message struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("Message");
            if (struct instanceof MessageEvent) {
                protocol.mo14136v3("messageEvent", 1, (byte) 12);
                com.x.dmv2.thriftjava.MessageEvent.ADAPTER.write(protocol, ((MessageEvent) struct).m76761getValue());
            } else if (struct instanceof MessageInstruction) {
                protocol.mo14136v3("messageInstruction", 2, (byte) 12);
                com.x.dmv2.thriftjava.MessageInstruction.ADAPTER.write(protocol, ((MessageInstruction) struct).m76762getValue());
            } else if (struct instanceof BatchedMessageEvents) {
                protocol.mo14136v3("batchedMessageEvents", 3, (byte) 12);
                com.x.dmv2.thriftjava.BatchedMessageEvents.ADAPTER.write(protocol, ((BatchedMessageEvents) struct).m76760getValue());
            } else if (!(struct instanceof Unknown)) {
                throw new NoWhenBranchMatchedException();
            }
            protocol.mo14134i0();
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/Message$MessageEvent;", "Lcom/x/dmv2/thriftjava/Message;", "value", "Lcom/x/dmv2/thriftjava/MessageEvent;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageEvent;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageEvent;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MessageEvent extends Message {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MessageEvent value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MessageEvent(@InterfaceC88464a com.x.dmv2.thriftjava.MessageEvent value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MessageEvent copy$default(MessageEvent messageEvent, com.x.dmv2.thriftjava.MessageEvent messageEvent2, int i, Object obj) {
            if ((i & 1) != 0) {
                messageEvent2 = messageEvent.value;
            }
            return messageEvent.copy(messageEvent2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MessageEvent getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MessageEvent copy(@InterfaceC88464a com.x.dmv2.thriftjava.MessageEvent value) {
            Intrinsics.m65272h(value, "value");
            return new MessageEvent(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MessageEvent) && Intrinsics.m65267c(this.value, ((MessageEvent) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MessageEvent m76761getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "Message(messageEvent=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/Message$MessageInstruction;", "Lcom/x/dmv2/thriftjava/Message;", "value", "Lcom/x/dmv2/thriftjava/MessageInstruction;", "<init>", "(Lcom/x/dmv2/thriftjava/MessageInstruction;)V", "getValue", "()Lcom/x/dmv2/thriftjava/MessageInstruction;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class MessageInstruction extends Message {

        @InterfaceC88464a
        private final com.x.dmv2.thriftjava.MessageInstruction value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public MessageInstruction(@InterfaceC88464a com.x.dmv2.thriftjava.MessageInstruction value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ MessageInstruction copy$default(MessageInstruction messageInstruction, com.x.dmv2.thriftjava.MessageInstruction messageInstruction2, int i, Object obj) {
            if ((i & 1) != 0) {
                messageInstruction2 = messageInstruction.value;
            }
            return messageInstruction.copy(messageInstruction2);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final com.x.dmv2.thriftjava.MessageInstruction getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final MessageInstruction copy(@InterfaceC88464a com.x.dmv2.thriftjava.MessageInstruction value) {
            Intrinsics.m65272h(value, "value");
            return new MessageInstruction(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof MessageInstruction) && Intrinsics.m65267c(this.value, ((MessageInstruction) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final com.x.dmv2.thriftjava.MessageInstruction m76762getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "Message(messageInstruction=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\n\u0002\u0010\u000e\n\u0000\bÆ\n\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0013\u0010\u0004\u001a\u00020\u00052\b\u0010\u0006\u001a\u0004\u0018\u00010\u0007HÖ\u0003J\t\u0010\b\u001a\u00020\tHÖ\u0001J\t\u0010\n\u001a\u00020\u000bHÖ\u0001¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/Message$Unknown;", "Lcom/x/dmv2/thriftjava/Message;", "<init>", "()V", "equals", "", "other", "", "hashCode", "", "toString", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Unknown extends Message {

        @InterfaceC88464a
        public static final Unknown INSTANCE = new Unknown();

        private Unknown() {
            super(null);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            return this == other || (other instanceof Unknown);
        }

        public int hashCode() {
            return 150203574;
        }

        @InterfaceC88464a
        public String toString() {
            return "Unknown";
        }
    }

    public /* synthetic */ Message(DefaultConstructorMarker defaultConstructorMarker) {
        this();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }

    private Message() {
    }
}
