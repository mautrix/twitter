package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
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

@Metadata(m64929d1 = {"\u0000:\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\t\n\u0000\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\t\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000 \u001e2\u00020\u0001:\u0002\u001f\u001eB\u001b\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004¢\u0006\u0004\b\u0006\u0010\u0007J\u0017\u0010\u000b\u001a\u00020\n2\u0006\u0010\t\u001a\u00020\bH\u0016¢\u0006\u0004\b\u000b\u0010\fJ\u0012\u0010\r\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\r\u0010\u000eJ\u0012\u0010\u000f\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u000f\u0010\u0010J(\u0010\u0011\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u0004HÆ\u0001¢\u0006\u0004\b\u0011\u0010\u0012J\u0010\u0010\u0013\u001a\u00020\u0004HÖ\u0001¢\u0006\u0004\b\u0013\u0010\u0010J\u0010\u0010\u0015\u001a\u00020\u0014HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u0016J\u001a\u0010\u001a\u001a\u00020\u00192\b\u0010\u0018\u001a\u0004\u0018\u00010\u0017HÖ\u0003¢\u0006\u0004\b\u001a\u0010\u001bR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001cR\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010\u001d¨\u0006 "}, m64930d2 = {"Lcom/x/dmv2/thriftjava/NicknameMessage;", "Lcom/bendb/thrifty/a;", "", "user_id", "", "nickname_text", "<init>", "(Ljava/lang/Long;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/Long;", "component2", "()Ljava/lang/String;", "copy", "(Ljava/lang/Long;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/NicknameMessage;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/Long;", "Ljava/lang/String;", "Companion", "NicknameMessageAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class NicknameMessage implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String nickname_text;

    @JvmField
    @InterfaceC88465b
    public final Long user_id;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new NicknameMessageAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/NicknameMessage$NicknameMessageAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/NicknameMessage;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/NicknameMessage;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/NicknameMessage;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class NicknameMessageAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public NicknameMessage m85665read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Long lValueOf = null;
            String string = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new NicknameMessage(lValueOf, string);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        C11272a.m14141a(protocol, b);
                    } else if (b == 11) {
                        string = protocol.readString();
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 10) {
                    lValueOf = Long.valueOf(protocol.mo14124H0());
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a NicknameMessage struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("NicknameMessage");
            if (struct.user_id != null) {
                protocol.mo14136v3("user_id", 1, (byte) 10);
                protocol.mo14121B3(struct.user_id.longValue());
            }
            if (struct.nickname_text != null) {
                protocol.mo14136v3("nickname_text", 2, (byte) 11);
                protocol.mo14137w0(struct.nickname_text);
            }
            protocol.mo14134i0();
        }
    }

    public NicknameMessage(@InterfaceC88465b Long l, @InterfaceC88465b String str) {
        this.user_id = l;
        this.nickname_text = str;
    }

    public static /* synthetic */ NicknameMessage copy$default(NicknameMessage nicknameMessage, Long l, String str, int i, Object obj) {
        if ((i & 1) != 0) {
            l = nicknameMessage.user_id;
        }
        if ((i & 2) != 0) {
            str = nicknameMessage.nickname_text;
        }
        return nicknameMessage.copy(l, str);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final Long getUser_id() {
        return this.user_id;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getNickname_text() {
        return this.nickname_text;
    }

    @InterfaceC88464a
    public final NicknameMessage copy(@InterfaceC88465b Long user_id, @InterfaceC88465b String nickname_text) {
        return new NicknameMessage(user_id, nickname_text);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof NicknameMessage)) {
            return false;
        }
        NicknameMessage nicknameMessage = (NicknameMessage) other;
        return Intrinsics.m65267c(this.user_id, nicknameMessage.user_id) && Intrinsics.m65267c(this.nickname_text, nicknameMessage.nickname_text);
    }

    public int hashCode() {
        Long l = this.user_id;
        int iHashCode = (l == null ? 0 : l.hashCode()) * 31;
        String str = this.nickname_text;
        return iHashCode + (str != null ? str.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "NicknameMessage(user_id=" + this.user_id + ", nickname_text=" + this.nickname_text + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}