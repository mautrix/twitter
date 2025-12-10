package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import android.gov.nist.javax.sip.C0024c;
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

@Metadata(m64929d1 = {"\u00004\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0007\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0006\b\u0086\b\u0018\u0000 \u00192\u00020\u0001:\u0002\u001a\u0019B\u0011\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002¢\u0006\u0004\b\u0004\u0010\u0005J\u0017\u0010\t\u001a\u00020\b2\u0006\u0010\u0007\u001a\u00020\u0006H\u0016¢\u0006\u0004\b\t\u0010\nJ\u0012\u0010\u000b\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000b\u0010\fJ\u001c\u0010\r\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\r\u0010\u000eJ\u0010\u0010\u000f\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u000f\u0010\fJ\u0010\u0010\u0011\u001a\u00020\u0010HÖ\u0001¢\u0006\u0004\b\u0011\u0010\u0012J\u001a\u0010\u0016\u001a\u00020\u00152\b\u0010\u0014\u001a\u0004\u0018\u00010\u0013HÖ\u0003¢\u0006\u0004\b\u0016\u0010\u0017R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u0018¨\u0006\u001b"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MarkConversationUnread;", "Lcom/bendb/thrifty/a;", "", "seen_until_sequence_id", "<init>", "(Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "copy", "(Ljava/lang/String;)Lcom/x/dmv2/thriftjava/MarkConversationUnread;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Companion", "MarkConversationUnreadAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class MarkConversationUnread implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String seen_until_sequence_id;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MarkConversationUnreadAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MarkConversationUnread$MarkConversationUnreadAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MarkConversationUnread;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MarkConversationUnread;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MarkConversationUnread;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MarkConversationUnreadAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MarkConversationUnread m85650read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new MarkConversationUnread(string);
                }
                if (c11265cMo14127V2.f38393b != 1) {
                    C11272a.m14141a(protocol, b);
                } else if (b == 11) {
                    string = protocol.readString();
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MarkConversationUnread struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MarkConversationUnread");
            if (struct.seen_until_sequence_id != null) {
                protocol.mo14136v3("seen_until_sequence_id", 1, (byte) 11);
                protocol.mo14137w0(struct.seen_until_sequence_id);
            }
            protocol.mo14134i0();
        }
    }

    public MarkConversationUnread(@InterfaceC88465b String str) {
        this.seen_until_sequence_id = str;
    }

    public static /* synthetic */ MarkConversationUnread copy$default(MarkConversationUnread markConversationUnread, String str, int i, Object obj) {
        if ((i & 1) != 0) {
            str = markConversationUnread.seen_until_sequence_id;
        }
        return markConversationUnread.copy(str);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getSeen_until_sequence_id() {
        return this.seen_until_sequence_id;
    }

    @InterfaceC88464a
    public final MarkConversationUnread copy(@InterfaceC88465b String seen_until_sequence_id) {
        return new MarkConversationUnread(seen_until_sequence_id);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        return (other instanceof MarkConversationUnread) && Intrinsics.m65267c(this.seen_until_sequence_id, ((MarkConversationUnread) other).seen_until_sequence_id);
    }

    public int hashCode() {
        String str = this.seen_until_sequence_id;
        if (str == null) {
            return 0;
        }
        return str.hashCode();
    }

    @InterfaceC88464a
    public String toString() {
        return C0024c.m31a("MarkConversationUnread(seen_until_sequence_id=", this.seen_until_sequence_id, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}