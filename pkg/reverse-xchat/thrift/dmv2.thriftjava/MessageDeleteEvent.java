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
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000>\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010 \n\u0002\u0010\u000e\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\n\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000  2\u00020\u0001:\u0002! B!\u0012\u000e\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002\u0012\b\u0010\u0006\u001a\u0004\u0018\u00010\u0005¢\u0006\u0004\b\u0007\u0010\bJ\u0017\u0010\f\u001a\u00020\u000b2\u0006\u0010\n\u001a\u00020\tH\u0016¢\u0006\u0004\b\f\u0010\rJ\u0018\u0010\u000e\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000e\u0010\u000fJ\u0012\u0010\u0010\u001a\u0004\u0018\u00010\u0005HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u0011J.\u0010\u0012\u001a\u00020\u00002\u0010\b\u0002\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00022\n\b\u0002\u0010\u0006\u001a\u0004\u0018\u00010\u0005HÆ\u0001¢\u0006\u0004\b\u0012\u0010\u0013J\u0010\u0010\u0014\u001a\u00020\u0003HÖ\u0001¢\u0006\u0004\b\u0014\u0010\u0015J\u0010\u0010\u0017\u001a\u00020\u0016HÖ\u0001¢\u0006\u0004\b\u0017\u0010\u0018J\u001a\u0010\u001c\u001a\u00020\u001b2\b\u0010\u001a\u001a\u0004\u0018\u00010\u0019HÖ\u0003¢\u0006\u0004\b\u001c\u0010\u001dR\u001c\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001eR\u0016\u0010\u0006\u001a\u0004\u0018\u00010\u00058\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010\u001f¨\u0006\""}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageDeleteEvent;", "Lcom/bendb/thrifty/a;", "", "", "sequence_ids", "Lcom/x/dmv2/thriftjava/DeleteMessageAction;", "delete_message_action", "<init>", "(Ljava/util/List;Lcom/x/dmv2/thriftjava/DeleteMessageAction;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/util/List;", "component2", "()Lcom/x/dmv2/thriftjava/DeleteMessageAction;", "copy", "(Ljava/util/List;Lcom/x/dmv2/thriftjava/DeleteMessageAction;)Lcom/x/dmv2/thriftjava/MessageDeleteEvent;", "toString", "()Ljava/lang/String;", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/util/List;", "Lcom/x/dmv2/thriftjava/DeleteMessageAction;", "Companion", "MessageDeleteEventAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class MessageDeleteEvent implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final DeleteMessageAction delete_message_action;

    @JvmField
    @InterfaceC88465b
    public final List sequence_ids;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MessageDeleteEventAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageDeleteEvent$MessageDeleteEventAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MessageDeleteEvent;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MessageDeleteEvent;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MessageDeleteEvent;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MessageDeleteEventAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MessageDeleteEvent m83743read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            ArrayList arrayList = null;
            DeleteMessageAction deleteMessageAction = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new MessageDeleteEvent(arrayList, deleteMessageAction);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        C11272a.m14141a(protocol, b);
                    } else if (b == 8) {
                        int iMo14132c4 = protocol.mo14132c4();
                        DeleteMessageAction deleteMessageActionFindByValue = DeleteMessageAction.INSTANCE.findByValue(iMo14132c4);
                        if (deleteMessageActionFindByValue == null) {
                            throw new ThriftException(ThriftException.EnumC11260b.PROTOCOL_ERROR, C0031b.m45c(iMo14132c4, "Unexpected value for enum type DeleteMessageAction: "));
                        }
                        deleteMessageAction = deleteMessageActionFindByValue;
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 15) {
                    int i = protocol.mo14130a2().f38395b;
                    ArrayList arrayList2 = new ArrayList(i);
                    for (int i2 = 0; i2 < i; i2++) {
                        arrayList2.add(protocol.readString());
                    }
                    arrayList = arrayList2;
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MessageDeleteEvent struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MessageDeleteEvent");
            if (struct.sequence_ids != null) {
                protocol.mo14136v3("sequence_ids", 1, (byte) 15);
                protocol.mo14128X0((byte) 11, struct.sequence_ids.size());
                Iterator it = struct.sequence_ids.iterator();
                while (it.hasNext()) {
                    protocol.mo14137w0((String) it.next());
                }
            }
            if (struct.delete_message_action != null) {
                protocol.mo14136v3("delete_message_action", 2, (byte) 8);
                protocol.mo14122C2(struct.delete_message_action.value);
            }
            protocol.mo14134i0();
        }
    }

    public MessageDeleteEvent(@InterfaceC88465b List list, @InterfaceC88465b DeleteMessageAction deleteMessageAction) {
        this.sequence_ids = list;
        this.delete_message_action = deleteMessageAction;
    }

    public static /* synthetic */ MessageDeleteEvent copy$default(MessageDeleteEvent messageDeleteEvent, List list, DeleteMessageAction deleteMessageAction, int i, Object obj) {
        if ((i & 1) != 0) {
            list = messageDeleteEvent.sequence_ids;
        }
        if ((i & 2) != 0) {
            deleteMessageAction = messageDeleteEvent.delete_message_action;
        }
        return messageDeleteEvent.copy(list, deleteMessageAction);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final List getSequence_ids() {
        return this.sequence_ids;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final DeleteMessageAction getDelete_message_action() {
        return this.delete_message_action;
    }

    @InterfaceC88464a
    public final MessageDeleteEvent copy(@InterfaceC88465b List sequence_ids, @InterfaceC88465b DeleteMessageAction delete_message_action) {
        return new MessageDeleteEvent(sequence_ids, delete_message_action);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof MessageDeleteEvent)) {
            return false;
        }
        MessageDeleteEvent messageDeleteEvent = (MessageDeleteEvent) other;
        return Intrinsics.m65267c(this.sequence_ids, messageDeleteEvent.sequence_ids) && this.delete_message_action == messageDeleteEvent.delete_message_action;
    }

    public int hashCode() {
        List list = this.sequence_ids;
        int iHashCode = (list == null ? 0 : list.hashCode()) * 31;
        DeleteMessageAction deleteMessageAction = this.delete_message_action;
        return iHashCode + (deleteMessageAction != null ? deleteMessageAction.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "MessageDeleteEvent(sequence_ids=" + this.sequence_ids + ", delete_message_action=" + this.delete_message_action + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
