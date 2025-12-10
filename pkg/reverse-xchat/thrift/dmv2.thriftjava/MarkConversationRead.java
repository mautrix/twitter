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

@Metadata(m64929d1 = {"\u0000:\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0000\n\u0002\u0010\t\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\t\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000 \u001e2\u00020\u0001:\u0002\u001f\u001eB\u001b\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004¢\u0006\u0004\b\u0006\u0010\u0007J\u0017\u0010\u000b\u001a\u00020\n2\u0006\u0010\t\u001a\u00020\bH\u0016¢\u0006\u0004\b\u000b\u0010\fJ\u0012\u0010\r\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\r\u0010\u000eJ\u0012\u0010\u000f\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u000f\u0010\u0010J(\u0010\u0011\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u0004HÆ\u0001¢\u0006\u0004\b\u0011\u0010\u0012J\u0010\u0010\u0013\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0013\u0010\u000eJ\u0010\u0010\u0015\u001a\u00020\u0014HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u0016J\u001a\u0010\u001a\u001a\u00020\u00192\b\u0010\u0018\u001a\u0004\u0018\u00010\u0017HÖ\u0003¢\u0006\u0004\b\u001a\u0010\u001bR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001cR\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010\u001d¨\u0006 "}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MarkConversationRead;", "Lcom/bendb/thrifty/a;", "", "seen_until_sequence_id", "", "seen_at_millis", "<init>", "(Ljava/lang/String;Ljava/lang/Long;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "()Ljava/lang/Long;", "copy", "(Ljava/lang/String;Ljava/lang/Long;)Lcom/x/dmv2/thriftjava/MarkConversationRead;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Ljava/lang/Long;", "Companion", "MarkConversationReadAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class MarkConversationRead implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final Long seen_at_millis;

    @JvmField
    @InterfaceC88465b
    public final String seen_until_sequence_id;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MarkConversationReadAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MarkConversationRead$MarkConversationReadAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MarkConversationRead;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MarkConversationRead;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MarkConversationRead;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MarkConversationReadAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MarkConversationRead m85648read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            Long lValueOf = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new MarkConversationRead(string, lValueOf);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        C11272a.m14141a(protocol, b);
                    } else if (b == 10) {
                        lValueOf = Long.valueOf(protocol.mo14124H0());
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

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MarkConversationRead struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MarkConversationRead");
            if (struct.seen_until_sequence_id != null) {
                protocol.mo14136v3("seen_until_sequence_id", 1, (byte) 11);
                protocol.mo14137w0(struct.seen_until_sequence_id);
            }
            if (struct.seen_at_millis != null) {
                protocol.mo14136v3("seen_at_millis", 2, (byte) 10);
                protocol.mo14121B3(struct.seen_at_millis.longValue());
            }
            protocol.mo14134i0();
        }
    }

    public MarkConversationRead(@InterfaceC88465b String str, @InterfaceC88465b Long l) {
        this.seen_until_sequence_id = str;
        this.seen_at_millis = l;
    }

    public static /* synthetic */ MarkConversationRead copy$default(MarkConversationRead markConversationRead, String str, Long l, int i, Object obj) {
        if ((i & 1) != 0) {
            str = markConversationRead.seen_until_sequence_id;
        }
        if ((i & 2) != 0) {
            l = markConversationRead.seen_at_millis;
        }
        return markConversationRead.copy(str, l);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getSeen_until_sequence_id() {
        return this.seen_until_sequence_id;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final Long getSeen_at_millis() {
        return this.seen_at_millis;
    }

    @InterfaceC88464a
    public final MarkConversationRead copy(@InterfaceC88465b String seen_until_sequence_id, @InterfaceC88465b Long seen_at_millis) {
        return new MarkConversationRead(seen_until_sequence_id, seen_at_millis);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof MarkConversationRead)) {
            return false;
        }
        MarkConversationRead markConversationRead = (MarkConversationRead) other;
        return Intrinsics.m65267c(this.seen_until_sequence_id, markConversationRead.seen_until_sequence_id) && Intrinsics.m65267c(this.seen_at_millis, markConversationRead.seen_at_millis);
    }

    public int hashCode() {
        String str = this.seen_until_sequence_id;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        Long l = this.seen_at_millis;
        return iHashCode + (l != null ? l.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "MarkConversationRead(seen_until_sequence_id=" + this.seen_until_sequence_id + ", seen_at_millis=" + this.seen_at_millis + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}