package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import com.x.composer.sensitivemedia.C68546u;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000D\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0000\n\u0002\u0010 \n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u000b\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\b\b\u0086\b\u0018\u0000 $2\u00020\u0001:\u0002%$B+\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\u000e\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u0004\u0012\b\u0010\b\u001a\u0004\u0018\u00010\u0007¢\u0006\u0004\b\t\u0010\nJ\u0017\u0010\u000e\u001a\u00020\r2\u0006\u0010\f\u001a\u00020\u000bH\u0016¢\u0006\u0004\b\u000e\u0010\u000fJ\u0012\u0010\u0010\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u0011J\u0018\u0010\u0012\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u0012\u0010\u0013J\u0012\u0010\u0014\u001a\u0004\u0018\u00010\u0007HÆ\u0003¢\u0006\u0004\b\u0014\u0010\u0015J:\u0010\u0016\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\u0010\b\u0002\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u00042\n\b\u0002\u0010\b\u001a\u0004\u0018\u00010\u0007HÆ\u0001¢\u0006\u0004\b\u0016\u0010\u0017J\u0010\u0010\u0018\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0018\u0010\u0011J\u0010\u0010\u001a\u001a\u00020\u0019HÖ\u0001¢\u0006\u0004\b\u001a\u0010\u001bJ\u001a\u0010\u001f\u001a\u00020\u001e2\b\u0010\u001d\u001a\u0004\u0018\u00010\u001cHÖ\u0003¢\u0006\u0004\b\u001f\u0010 R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010!R\u001c\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010\"R\u0016\u0010\b\u001a\u0004\u0018\u00010\u00078\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\b\u0010#¨\u0006&"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationKeyChangeEvent;", "Lcom/bendb/thrifty/a;", "", "conversation_key_version", "", "Lcom/x/dmv2/thriftjava/ConversationParticipantKey;", "conversation_participant_keys", "Lcom/x/dmv2/thriftjava/KeyRotation;", "ratchet_tree", "<init>", "(Ljava/lang/String;Ljava/util/List;Lcom/x/dmv2/thriftjava/KeyRotation;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "()Ljava/util/List;", "component3", "()Lcom/x/dmv2/thriftjava/KeyRotation;", "copy", "(Ljava/lang/String;Ljava/util/List;Lcom/x/dmv2/thriftjava/KeyRotation;)Lcom/x/dmv2/thriftjava/ConversationKeyChangeEvent;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Ljava/util/List;", "Lcom/x/dmv2/thriftjava/KeyRotation;", "Companion", "ConversationKeyChangeEventAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class ConversationKeyChangeEvent implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String conversation_key_version;

    @JvmField
    @InterfaceC88465b
    public final List conversation_participant_keys;

    @JvmField
    @InterfaceC88465b
    public final KeyRotation ratchet_tree;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new ConversationKeyChangeEventAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationKeyChangeEvent$ConversationKeyChangeEventAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/ConversationKeyChangeEvent;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/ConversationKeyChangeEvent;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/ConversationKeyChangeEvent;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class ConversationKeyChangeEventAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public ConversationKeyChangeEvent m85581read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            ArrayList arrayList = null;
            KeyRotation keyRotation = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new ConversationKeyChangeEvent(string, arrayList, keyRotation);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            C11272a.m14141a(protocol, b);
                        } else if (b == 12) {
                            keyRotation = (KeyRotation) KeyRotation.ADAPTER.read(protocol);
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 15) {
                        int i = protocol.mo14130a2().f38395b;
                        ArrayList arrayList2 = new ArrayList(i);
                        for (int i2 = 0; i2 < i; i2++) {
                            arrayList2.add((ConversationParticipantKey) ConversationParticipantKey.ADAPTER.read(protocol));
                        }
                        arrayList = arrayList2;
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 11) {
                    string = protocol.readString();
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a ConversationKeyChangeEvent struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("ConversationKeyChangeEvent");
            if (struct.conversation_key_version != null) {
                protocol.mo14136v3("conversation_key_version", 1, (byte) 11);
                protocol.mo14137w0(struct.conversation_key_version);
            }
            if (struct.conversation_participant_keys != null) {
                protocol.mo14136v3("conversation_participant_keys", 2, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.conversation_participant_keys.size());
                Iterator it = struct.conversation_participant_keys.iterator();
                while (it.hasNext()) {
                    ConversationParticipantKey.ADAPTER.write(protocol, (ConversationParticipantKey) it.next());
                }
            }
            if (struct.ratchet_tree != null) {
                protocol.mo14136v3("ratchet_tree", 3, (byte) 12);
                KeyRotation.ADAPTER.write(protocol, struct.ratchet_tree);
            }
            protocol.mo14134i0();
        }
    }

    public ConversationKeyChangeEvent(@InterfaceC88465b String str, @InterfaceC88465b List list, @InterfaceC88465b KeyRotation keyRotation) {
        this.conversation_key_version = str;
        this.conversation_participant_keys = list;
        this.ratchet_tree = keyRotation;
    }

    public static /* synthetic */ ConversationKeyChangeEvent copy$default(ConversationKeyChangeEvent conversationKeyChangeEvent, String str, List list, KeyRotation keyRotation, int i, Object obj) {
        if ((i & 1) != 0) {
            str = conversationKeyChangeEvent.conversation_key_version;
        }
        if ((i & 2) != 0) {
            list = conversationKeyChangeEvent.conversation_participant_keys;
        }
        if ((i & 4) != 0) {
            keyRotation = conversationKeyChangeEvent.ratchet_tree;
        }
        return conversationKeyChangeEvent.copy(str, list, keyRotation);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getConversation_key_version() {
        return this.conversation_key_version;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final List getConversation_participant_keys() {
        return this.conversation_participant_keys;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final KeyRotation getRatchet_tree() {
        return this.ratchet_tree;
    }

    @InterfaceC88464a
    public final ConversationKeyChangeEvent copy(@InterfaceC88465b String conversation_key_version, @InterfaceC88465b List conversation_participant_keys, @InterfaceC88465b KeyRotation ratchet_tree) {
        return new ConversationKeyChangeEvent(conversation_key_version, conversation_participant_keys, ratchet_tree);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof ConversationKeyChangeEvent)) {
            return false;
        }
        ConversationKeyChangeEvent conversationKeyChangeEvent = (ConversationKeyChangeEvent) other;
        return Intrinsics.m65267c(this.conversation_key_version, conversationKeyChangeEvent.conversation_key_version) && Intrinsics.m65267c(this.conversation_participant_keys, conversationKeyChangeEvent.conversation_participant_keys) && Intrinsics.m65267c(this.ratchet_tree, conversationKeyChangeEvent.ratchet_tree);
    }

    public int hashCode() {
        String str = this.conversation_key_version;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        List list = this.conversation_participant_keys;
        int iHashCode2 = (iHashCode + (list == null ? 0 : list.hashCode())) * 31;
        KeyRotation keyRotation = this.ratchet_tree;
        return iHashCode2 + (keyRotation != null ? keyRotation.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        String str = this.conversation_key_version;
        List list = this.conversation_participant_keys;
        KeyRotation keyRotation = this.ratchet_tree;
        StringBuilder sbM56620a = C68546u.m56620a("ConversationKeyChangeEvent(conversation_key_version=", str, ", conversation_participant_keys=", list, ", ratchet_tree=");
        sbM56620a.append(keyRotation);
        sbM56620a.append(Separators.RPAREN);
        return sbM56620a.toString();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
