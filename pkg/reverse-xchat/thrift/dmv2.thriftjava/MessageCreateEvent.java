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
import okio.C87081h;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000R\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u000e\n\u0000\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\t\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010 \n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0013\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0002\b\f\b\u0086\b\u0018\u0000 62\u00020\u0001:\u000276B]\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004\u0012\b\u0010\u0007\u001a\u0004\u0018\u00010\u0006\u0012\b\u0010\t\u001a\u0004\u0018\u00010\b\u0012\b\u0010\n\u001a\u0004\u0018\u00010\b\u0012\b\u0010\u000b\u001a\u0004\u0018\u00010\u0006\u0012\b\u0010\r\u001a\u0004\u0018\u00010\f\u0012\u000e\u0010\u0010\u001a\n\u0012\u0004\u0012\u00020\u000f\u0018\u00010\u000e¢\u0006\u0004\b\u0011\u0010\u0012J\u0017\u0010\u0016\u001a\u00020\u00152\u0006\u0010\u0014\u001a\u00020\u0013H\u0016¢\u0006\u0004\b\u0016\u0010\u0017J\u0012\u0010\u0018\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0018\u0010\u0019J\u0012\u0010\u001a\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u001a\u0010\u001bJ\u0012\u0010\u001c\u001a\u0004\u0018\u00010\u0006HÆ\u0003¢\u0006\u0004\b\u001c\u0010\u001dJ\u0012\u0010\u001e\u001a\u0004\u0018\u00010\bHÆ\u0003¢\u0006\u0004\b\u001e\u0010\u001fJ\u0012\u0010 \u001a\u0004\u0018\u00010\bHÆ\u0003¢\u0006\u0004\b \u0010\u001fJ\u0012\u0010!\u001a\u0004\u0018\u00010\u0006HÆ\u0003¢\u0006\u0004\b!\u0010\u001dJ\u0012\u0010\"\u001a\u0004\u0018\u00010\fHÆ\u0003¢\u0006\u0004\b\"\u0010#J\u0018\u0010$\u001a\n\u0012\u0004\u0012\u00020\u000f\u0018\u00010\u000eHÆ\u0003¢\u0006\u0004\b$\u0010%Jv\u0010&\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u00042\n\b\u0002\u0010\u0007\u001a\u0004\u0018\u00010\u00062\n\b\u0002\u0010\t\u001a\u0004\u0018\u00010\b2\n\b\u0002\u0010\n\u001a\u0004\u0018\u00010\b2\n\b\u0002\u0010\u000b\u001a\u0004\u0018\u00010\u00062\n\b\u0002\u0010\r\u001a\u0004\u0018\u00010\f2\u0010\b\u0002\u0010\u0010\u001a\n\u0012\u0004\u0012\u00020\u000f\u0018\u00010\u000eHÆ\u0001¢\u0006\u0004\b&\u0010'J\u0010\u0010(\u001a\u00020\u0004HÖ\u0001¢\u0006\u0004\b(\u0010\u001bJ\u0010\u0010*\u001a\u00020)HÖ\u0001¢\u0006\u0004\b*\u0010+J\u001a\u0010.\u001a\u00020\u00062\b\u0010-\u001a\u0004\u0018\u00010,HÖ\u0003¢\u0006\u0004\b.\u0010/R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u00100R\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u00101R\u0016\u0010\u0007\u001a\u0004\u0018\u00010\u00068\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0007\u00102R\u0016\u0010\t\u001a\u0004\u0018\u00010\b8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\t\u00103R\u0016\u0010\n\u001a\u0004\u0018\u00010\b8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\n\u00103R\u0016\u0010\u000b\u001a\u0004\u0018\u00010\u00068\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u000b\u00102R\u0016\u0010\r\u001a\u0004\u0018\u00010\f8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\r\u00104R\u001c\u0010\u0010\u001a\n\u0012\u0004\u0012\u00020\u000f\u0018\u00010\u000e8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0010\u00105¨\u00068"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageCreateEvent;", "Lcom/bendb/thrifty/a;", "Lokio/h;", "contents", "", "conversation_key_version", "", "should_notify", "", "ttl_msec", "delivered_at_msec", "is_pending_public_key", "Lcom/x/dmv2/thriftjava/EventQueuePriority;", "priority", "", "Lcom/x/dmv2/thriftjava/AdditionalAction;", "additional_action_list", "<init>", "(Lokio/h;Ljava/lang/String;Ljava/lang/Boolean;Ljava/lang/Long;Ljava/lang/Long;Ljava/lang/Boolean;Lcom/x/dmv2/thriftjava/EventQueuePriority;Ljava/util/List;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Lokio/h;", "component2", "()Ljava/lang/String;", "component3", "()Ljava/lang/Boolean;", "component4", "()Ljava/lang/Long;", "component5", "component6", "component7", "()Lcom/x/dmv2/thriftjava/EventQueuePriority;", "component8", "()Ljava/util/List;", "copy", "(Lokio/h;Ljava/lang/String;Ljava/lang/Boolean;Ljava/lang/Long;Ljava/lang/Long;Ljava/lang/Boolean;Lcom/x/dmv2/thriftjava/EventQueuePriority;Ljava/util/List;)Lcom/x/dmv2/thriftjava/MessageCreateEvent;", "toString", "", "hashCode", "()I", "", "other", "equals", "(Ljava/lang/Object;)Z", "Lokio/h;", "Ljava/lang/String;", "Ljava/lang/Boolean;", "Ljava/lang/Long;", "Lcom/x/dmv2/thriftjava/EventQueuePriority;", "Ljava/util/List;", "Companion", "MessageCreateEventAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class MessageCreateEvent implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final List additional_action_list;

    @JvmField
    @InterfaceC88465b
    public final C87081h contents;

    @JvmField
    @InterfaceC88465b
    public final String conversation_key_version;

    @JvmField
    @InterfaceC88465b
    public final Long delivered_at_msec;

    @JvmField
    @InterfaceC88465b
    public final Boolean is_pending_public_key;

    @JvmField
    @InterfaceC88465b
    public final EventQueuePriority priority;

    @JvmField
    @InterfaceC88465b
    public final Boolean should_notify;

    @JvmField
    @InterfaceC88465b
    public final Long ttl_msec;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MessageCreateEventAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageCreateEvent$MessageCreateEventAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MessageCreateEvent;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MessageCreateEvent;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MessageCreateEvent;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MessageCreateEventAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MessageCreateEvent m85654read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            C87081h c87081hMo14123G1 = null;
            String string = null;
            Boolean boolValueOf = null;
            Long lValueOf = null;
            Long lValueOf2 = null;
            Boolean boolValueOf2 = null;
            EventQueuePriority eventQueuePriority = null;
            ArrayList arrayList = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new MessageCreateEvent(c87081hMo14123G1, string, boolValueOf, lValueOf, lValueOf2, boolValueOf2, eventQueuePriority, arrayList);
                }
                switch (c11265cMo14127V2.f38393b) {
                    case 100:
                        if (b == 11) {
                            c87081hMo14123G1 = protocol.mo14123G1();
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                            break;
                        }
                    case 101:
                        if (b == 11) {
                            string = protocol.readString();
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                            break;
                        }
                    case 102:
                        if (b == 2) {
                            boolValueOf = Boolean.valueOf(protocol.readBool());
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                            break;
                        }
                    case 103:
                        if (b == 10) {
                            lValueOf = Long.valueOf(protocol.mo14124H0());
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                            break;
                        }
                    case 104:
                        if (b == 10) {
                            lValueOf2 = Long.valueOf(protocol.mo14124H0());
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                            break;
                        }
                    case 105:
                        if (b == 2) {
                            boolValueOf2 = Boolean.valueOf(protocol.readBool());
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                            break;
                        }
                    case 106:
                        if (b == 8) {
                            int iMo14132c4 = protocol.mo14132c4();
                            EventQueuePriority eventQueuePriorityFindByValue = EventQueuePriority.INSTANCE.findByValue(iMo14132c4);
                            if (eventQueuePriorityFindByValue == null) {
                                throw new ThriftException(ThriftException.EnumC11260b.PROTOCOL_ERROR, C0031b.m45c(iMo14132c4, "Unexpected value for enum type EventQueuePriority: "));
                            }
                            eventQueuePriority = eventQueuePriorityFindByValue;
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                            break;
                        }
                    case 107:
                        if (b == 15) {
                            int i = protocol.mo14130a2().f38395b;
                            ArrayList arrayList2 = new ArrayList(i);
                            for (int i2 = 0; i2 < i; i2++) {
                                int iMo14132c42 = protocol.mo14132c4();
                                AdditionalAction additionalActionFindByValue = AdditionalAction.INSTANCE.findByValue(iMo14132c42);
                                if (additionalActionFindByValue == null) {
                                    throw new ThriftException(ThriftException.EnumC11260b.PROTOCOL_ERROR, C0031b.m45c(iMo14132c42, "Unexpected value for enum type AdditionalAction: "));
                                }
                                arrayList2.add(additionalActionFindByValue);
                            }
                            arrayList = arrayList2;
                            break;
                        } else {
                            C11272a.m14141a(protocol, b);
                            break;
                        }
                    default:
                        C11272a.m14141a(protocol, b);
                        break;
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MessageCreateEvent struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MessageCreateEvent");
            if (struct.contents != null) {
                protocol.mo14136v3("contents", 100, (byte) 11);
                protocol.mo14138y0(struct.contents);
            }
            if (struct.conversation_key_version != null) {
                protocol.mo14136v3("conversation_key_version", 101, (byte) 11);
                protocol.mo14137w0(struct.conversation_key_version);
            }
            if (struct.should_notify != null) {
                protocol.mo14136v3("should_notify", 102, (byte) 2);
                protocol.mo14125P1(struct.should_notify.booleanValue());
            }
            if (struct.ttl_msec != null) {
                protocol.mo14136v3("ttl_msec", 103, (byte) 10);
                protocol.mo14121B3(struct.ttl_msec.longValue());
            }
            if (struct.delivered_at_msec != null) {
                protocol.mo14136v3("delivered_at_msec", 104, (byte) 10);
                protocol.mo14121B3(struct.delivered_at_msec.longValue());
            }
            if (struct.is_pending_public_key != null) {
                protocol.mo14136v3("is_pending_public_key", 105, (byte) 2);
                protocol.mo14125P1(struct.is_pending_public_key.booleanValue());
            }
            if (struct.priority != null) {
                protocol.mo14136v3("priority", 106, (byte) 8);
                protocol.mo14122C2(struct.priority.value);
            }
            if (struct.additional_action_list != null) {
                protocol.mo14136v3("additional_action_list", 107, (byte) 15);
                protocol.mo14128X0((byte) 8, struct.additional_action_list.size());
                Iterator it = struct.additional_action_list.iterator();
                while (it.hasNext()) {
                    protocol.mo14122C2(((AdditionalAction) it.next()).value);
                }
            }
            protocol.mo14134i0();
        }
    }

    public MessageCreateEvent(@InterfaceC88465b C87081h c87081h, @InterfaceC88465b String str, @InterfaceC88465b Boolean bool, @InterfaceC88465b Long l, @InterfaceC88465b Long l2, @InterfaceC88465b Boolean bool2, @InterfaceC88465b EventQueuePriority eventQueuePriority, @InterfaceC88465b List list) {
        this.contents = c87081h;
        this.conversation_key_version = str;
        this.should_notify = bool;
        this.ttl_msec = l;
        this.delivered_at_msec = l2;
        this.is_pending_public_key = bool2;
        this.priority = eventQueuePriority;
        this.additional_action_list = list;
    }

    public static /* synthetic */ MessageCreateEvent copy$default(MessageCreateEvent messageCreateEvent, C87081h c87081h, String str, Boolean bool, Long l, Long l2, Boolean bool2, EventQueuePriority eventQueuePriority, List list, int i, Object obj) {
        return messageCreateEvent.copy((i & 1) != 0 ? messageCreateEvent.contents : c87081h, (i & 2) != 0 ? messageCreateEvent.conversation_key_version : str, (i & 4) != 0 ? messageCreateEvent.should_notify : bool, (i & 8) != 0 ? messageCreateEvent.ttl_msec : l, (i & 16) != 0 ? messageCreateEvent.delivered_at_msec : l2, (i & 32) != 0 ? messageCreateEvent.is_pending_public_key : bool2, (i & 64) != 0 ? messageCreateEvent.priority : eventQueuePriority, (i & 128) != 0 ? messageCreateEvent.additional_action_list : list);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final C87081h getContents() {
        return this.contents;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getConversation_key_version() {
        return this.conversation_key_version;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final Boolean getShould_notify() {
        return this.should_notify;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final Long getTtl_msec() {
        return this.ttl_msec;
    }

    @InterfaceC88465b
    /* renamed from: component5, reason: from getter */
    public final Long getDelivered_at_msec() {
        return this.delivered_at_msec;
    }

    @InterfaceC88465b
    /* renamed from: component6, reason: from getter */
    public final Boolean getIs_pending_public_key() {
        return this.is_pending_public_key;
    }

    @InterfaceC88465b
    /* renamed from: component7, reason: from getter */
    public final EventQueuePriority getPriority() {
        return this.priority;
    }

    @InterfaceC88465b
    /* renamed from: component8, reason: from getter */
    public final List getAdditional_action_list() {
        return this.additional_action_list;
    }

    @InterfaceC88464a
    public final MessageCreateEvent copy(@InterfaceC88465b C87081h contents, @InterfaceC88465b String conversation_key_version, @InterfaceC88465b Boolean should_notify, @InterfaceC88465b Long ttl_msec, @InterfaceC88465b Long delivered_at_msec, @InterfaceC88465b Boolean is_pending_public_key, @InterfaceC88465b EventQueuePriority priority, @InterfaceC88465b List additional_action_list) {
        return new MessageCreateEvent(contents, conversation_key_version, should_notify, ttl_msec, delivered_at_msec, is_pending_public_key, priority, additional_action_list);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof MessageCreateEvent)) {
            return false;
        }
        MessageCreateEvent messageCreateEvent = (MessageCreateEvent) other;
        return Intrinsics.m65267c(this.contents, messageCreateEvent.contents) && Intrinsics.m65267c(this.conversation_key_version, messageCreateEvent.conversation_key_version) && Intrinsics.m65267c(this.should_notify, messageCreateEvent.should_notify) && Intrinsics.m65267c(this.ttl_msec, messageCreateEvent.ttl_msec) && Intrinsics.m65267c(this.delivered_at_msec, messageCreateEvent.delivered_at_msec) && Intrinsics.m65267c(this.is_pending_public_key, messageCreateEvent.is_pending_public_key) && this.priority == messageCreateEvent.priority && Intrinsics.m65267c(this.additional_action_list, messageCreateEvent.additional_action_list);
    }

    public int hashCode() {
        C87081h c87081h = this.contents;
        int iHashCode = (c87081h == null ? 0 : c87081h.hashCode()) * 31;
        String str = this.conversation_key_version;
        int iHashCode2 = (iHashCode + (str == null ? 0 : str.hashCode())) * 31;
        Boolean bool = this.should_notify;
        int iHashCode3 = (iHashCode2 + (bool == null ? 0 : bool.hashCode())) * 31;
        Long l = this.ttl_msec;
        int iHashCode4 = (iHashCode3 + (l == null ? 0 : l.hashCode())) * 31;
        Long l2 = this.delivered_at_msec;
        int iHashCode5 = (iHashCode4 + (l2 == null ? 0 : l2.hashCode())) * 31;
        Boolean bool2 = this.is_pending_public_key;
        int iHashCode6 = (iHashCode5 + (bool2 == null ? 0 : bool2.hashCode())) * 31;
        EventQueuePriority eventQueuePriority = this.priority;
        int iHashCode7 = (iHashCode6 + (eventQueuePriority == null ? 0 : eventQueuePriority.hashCode())) * 31;
        List list = this.additional_action_list;
        return iHashCode7 + (list != null ? list.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "MessageCreateEvent(contents=" + this.contents + ", conversation_key_version=" + this.conversation_key_version + ", should_notify=" + this.should_notify + ", ttl_msec=" + this.ttl_msec + ", delivered_at_msec=" + this.delivered_at_msec + ", is_pending_public_key=" + this.is_pending_public_key + ", priority=" + this.priority + ", additional_action_list=" + this.additional_action_list + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
