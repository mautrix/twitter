package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import android.gov.nist.javax.sip.header.C0031b;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.ThriftException;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000<\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0006\n\u0002\u0010\u000e\n\u0002\b\u0002\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0006\b\u0086\b\u0018\u0000 \u001b2\u00020\u0001:\u0002\u001c\u001bB\u0011\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002¢\u0006\u0004\b\u0004\u0010\u0005J\u0017\u0010\t\u001a\u00020\b2\u0006\u0010\u0007\u001a\u00020\u0006H\u0016¢\u0006\u0004\b\t\u0010\nJ\u0012\u0010\u000b\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000b\u0010\fJ\u001c\u0010\r\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\r\u0010\u000eJ\u0010\u0010\u0010\u001a\u00020\u000fHÖ\u0001¢\u0006\u0004\b\u0010\u0010\u0011J\u0010\u0010\u0013\u001a\u00020\u0012HÖ\u0001¢\u0006\u0004\b\u0013\u0010\u0014J\u001a\u0010\u0018\u001a\u00020\u00172\b\u0010\u0016\u001a\u0004\u0018\u00010\u0015HÖ\u0003¢\u0006\u0004\b\u0018\u0010\u0019R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001a¨\u0006\u001d"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageFailureEvent;", "Lcom/bendb/thrifty/a;", "Lcom/x/dmv2/thriftjava/FailureType;", "failure_type", "<init>", "(Lcom/x/dmv2/thriftjava/FailureType;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Lcom/x/dmv2/thriftjava/FailureType;", "copy", "(Lcom/x/dmv2/thriftjava/FailureType;)Lcom/x/dmv2/thriftjava/MessageFailureEvent;", "", "toString", "()Ljava/lang/String;", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Lcom/x/dmv2/thriftjava/FailureType;", "Companion", "MessageFailureEventAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class MessageFailureEvent implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final FailureType failure_type;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MessageFailureEventAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageFailureEvent$MessageFailureEventAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MessageFailureEvent;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MessageFailureEvent;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MessageFailureEvent;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MessageFailureEventAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MessageFailureEvent m83752read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            FailureType failureType = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new MessageFailureEvent(failureType);
                }
                if (c11265cMo14127V2.f38393b != 1) {
                    C11272a.m14141a(protocol, b);
                } else if (b == 8) {
                    int iMo14132c4 = protocol.mo14132c4();
                    FailureType failureTypeFindByValue = FailureType.INSTANCE.findByValue(iMo14132c4);
                    if (failureTypeFindByValue == null) {
                        throw new ThriftException(ThriftException.EnumC11260b.PROTOCOL_ERROR, C0031b.m45c(iMo14132c4, "Unexpected value for enum type FailureType: "));
                    }
                    failureType = failureTypeFindByValue;
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MessageFailureEvent struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MessageFailureEvent");
            if (struct.failure_type != null) {
                protocol.mo14136v3("failure_type", 1, (byte) 8);
                protocol.mo14122C2(struct.failure_type.value);
            }
            protocol.mo14134i0();
        }
    }

    public MessageFailureEvent(@InterfaceC88465b FailureType failureType) {
        this.failure_type = failureType;
    }

    public static /* synthetic */ MessageFailureEvent copy$default(MessageFailureEvent messageFailureEvent, FailureType failureType, int i, Object obj) {
        if ((i & 1) != 0) {
            failureType = messageFailureEvent.failure_type;
        }
        return messageFailureEvent.copy(failureType);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final FailureType getFailure_type() {
        return this.failure_type;
    }

    @InterfaceC88464a
    public final MessageFailureEvent copy(@InterfaceC88465b FailureType failure_type) {
        return new MessageFailureEvent(failure_type);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        return (other instanceof MessageFailureEvent) && this.failure_type == ((MessageFailureEvent) other).failure_type;
    }

    public int hashCode() {
        FailureType failureType = this.failure_type;
        if (failureType == null) {
            return 0;
        }
        return failureType.hashCode();
    }

    @InterfaceC88464a
    public String toString() {
        return "MessageFailureEvent(failure_type=" + this.failure_type + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}