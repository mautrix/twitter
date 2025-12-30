package com.x.dmv2.thriftjava;

import android.gov.nist.core.C0003b;
import android.gov.nist.core.Separators;
import android.gov.nist.core.net.C0009a;
import com.bendb.thrifty.InterfaceC11261a;
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

@Metadata(m64929d1 = {"\u00004\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0002\b\u0005\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\t\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0006\b\u0086\b\u0018\u0000 \u001d2\u00020\u0001:\u0002\u001e\u001dB%\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0004\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0002¢\u0006\u0004\b\u0006\u0010\u0007J\u0017\u0010\u000b\u001a\u00020\n2\u0006\u0010\t\u001a\u00020\bH\u0016¢\u0006\u0004\b\u000b\u0010\fJ\u0012\u0010\r\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\r\u0010\u000eJ\u0012\u0010\u000f\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000f\u0010\u000eJ\u0012\u0010\u0010\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u000eJ4\u0010\u0011\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0004\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\u0011\u0010\u0012J\u0010\u0010\u0013\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0013\u0010\u000eJ\u0010\u0010\u0015\u001a\u00020\u0014HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u0016J\u001a\u0010\u001a\u001a\u00020\u00192\b\u0010\u0018\u001a\u0004\u0018\u00010\u0017HÖ\u0003¢\u0006\u0004\b\u001a\u0010\u001bR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001cR\u0016\u0010\u0004\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001cR\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010\u001c¨\u0006\u001f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationParticipantKey;", "Lcom/bendb/thrifty/a;", "", "user_id", "encrypted_conversation_key", "public_key_version", "<init>", "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "component3", "copy", "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/ConversationParticipantKey;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Companion", "ConversationParticipantKeyAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class ConversationParticipantKey implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String encrypted_conversation_key;

    @JvmField
    @InterfaceC88465b
    public final String public_key_version;

    @JvmField
    @InterfaceC88465b
    public final String user_id;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new ConversationParticipantKeyAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ConversationParticipantKey$ConversationParticipantKeyAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/ConversationParticipantKey;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/ConversationParticipantKey;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/ConversationParticipantKey;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class ConversationParticipantKeyAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public ConversationParticipantKey m85641read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            String string2 = null;
            String string3 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new ConversationParticipantKey(string, string2, string3);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            C11272a.m14141a(protocol, b);
                        } else if (b == 11) {
                            string3 = protocol.readString();
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 11) {
                        string2 = protocol.readString();
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

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a ConversationParticipantKey struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("ConversationParticipantKey");
            if (struct.user_id != null) {
                protocol.mo14136v3("user_id", 1, (byte) 11);
                protocol.mo14137w0(struct.user_id);
            }
            if (struct.encrypted_conversation_key != null) {
                protocol.mo14136v3("encrypted_conversation_key", 2, (byte) 11);
                protocol.mo14137w0(struct.encrypted_conversation_key);
            }
            if (struct.public_key_version != null) {
                protocol.mo14136v3("public_key_version", 3, (byte) 11);
                protocol.mo14137w0(struct.public_key_version);
            }
            protocol.mo14134i0();
        }
    }

    public ConversationParticipantKey(@InterfaceC88465b String str, @InterfaceC88465b String str2, @InterfaceC88465b String str3) {
        this.user_id = str;
        this.encrypted_conversation_key = str2;
        this.public_key_version = str3;
    }

    public static /* synthetic */ ConversationParticipantKey copy$default(ConversationParticipantKey conversationParticipantKey, String str, String str2, String str3, int i, Object obj) {
        if ((i & 1) != 0) {
            str = conversationParticipantKey.user_id;
        }
        if ((i & 2) != 0) {
            str2 = conversationParticipantKey.encrypted_conversation_key;
        }
        if ((i & 4) != 0) {
            str3 = conversationParticipantKey.public_key_version;
        }
        return conversationParticipantKey.copy(str, str2, str3);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getUser_id() {
        return this.user_id;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getEncrypted_conversation_key() {
        return this.encrypted_conversation_key;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final String getPublic_key_version() {
        return this.public_key_version;
    }

    @InterfaceC88464a
    public final ConversationParticipantKey copy(@InterfaceC88465b String user_id, @InterfaceC88465b String encrypted_conversation_key, @InterfaceC88465b String public_key_version) {
        return new ConversationParticipantKey(user_id, encrypted_conversation_key, public_key_version);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof ConversationParticipantKey)) {
            return false;
        }
        ConversationParticipantKey conversationParticipantKey = (ConversationParticipantKey) other;
        return Intrinsics.m65267c(this.user_id, conversationParticipantKey.user_id) && Intrinsics.m65267c(this.encrypted_conversation_key, conversationParticipantKey.encrypted_conversation_key) && Intrinsics.m65267c(this.public_key_version, conversationParticipantKey.public_key_version);
    }

    public int hashCode() {
        String str = this.user_id;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        String str2 = this.encrypted_conversation_key;
        int iHashCode2 = (iHashCode + (str2 == null ? 0 : str2.hashCode())) * 31;
        String str3 = this.public_key_version;
        return iHashCode2 + (str3 != null ? str3.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        String str = this.user_id;
        String str2 = this.encrypted_conversation_key;
        return C0003b.m4b(C0009a.m11b("ConversationParticipantKey(user_id=", str, ", encrypted_conversation_key=", str2, ", public_key_version="), this.public_key_version, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
