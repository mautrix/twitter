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

@Metadata(m64929d1 = {"\u0000<\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\t\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\b\n\u0002\u0010\u000e\n\u0002\b\u0002\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0002\b\b\b\u0086\b\u0018\u0000 \u001f2\u00020\u0001:\u0002 \u001fB\u001b\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004¢\u0006\u0004\b\u0006\u0010\u0007J\u0017\u0010\u000b\u001a\u00020\n2\u0006\u0010\t\u001a\u00020\bH\u0016¢\u0006\u0004\b\u000b\u0010\fJ\u0012\u0010\r\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\r\u0010\u000eJ\u0012\u0010\u000f\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u000f\u0010\u0010J(\u0010\u0011\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u0004HÆ\u0001¢\u0006\u0004\b\u0011\u0010\u0012J\u0010\u0010\u0014\u001a\u00020\u0013HÖ\u0001¢\u0006\u0004\b\u0014\u0010\u0015J\u0010\u0010\u0017\u001a\u00020\u0016HÖ\u0001¢\u0006\u0004\b\u0017\u0010\u0018J\u001a\u0010\u001b\u001a\u00020\u00042\b\u0010\u001a\u001a\u0004\u0018\u00010\u0019HÖ\u0003¢\u0006\u0004\b\u001b\u0010\u001cR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001dR\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010\u001e¨\u0006!"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/AVCallMissed;", "Lcom/bendb/thrifty/a;", "", "sent_at_millis", "", "is_audio_only", "<init>", "(Ljava/lang/Long;Ljava/lang/Boolean;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/Long;", "component2", "()Ljava/lang/Boolean;", "copy", "(Ljava/lang/Long;Ljava/lang/Boolean;)Lcom/x/dmv2/thriftjava/AVCallMissed;", "", "toString", "()Ljava/lang/String;", "", "hashCode", "()I", "", "other", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/Long;", "Ljava/lang/Boolean;", "Companion", "AVCallMissedAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class AVCallMissed implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final Boolean is_audio_only;

    @JvmField
    @InterfaceC88465b
    public final Long sent_at_millis;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new AVCallMissedAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/AVCallMissed$AVCallMissedAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/AVCallMissed;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/AVCallMissed;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/AVCallMissed;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class AVCallMissedAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public AVCallMissed m85637read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Long lValueOf = null;
            Boolean boolValueOf = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new AVCallMissed(lValueOf, boolValueOf);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        C11272a.m14141a(protocol, b);
                    } else if (b == 2) {
                        boolValueOf = Boolean.valueOf(protocol.readBool());
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

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a AVCallMissed struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("AVCallMissed");
            if (struct.sent_at_millis != null) {
                protocol.mo14136v3("sent_at_millis", 1, (byte) 10);
                protocol.mo14121B3(struct.sent_at_millis.longValue());
            }
            if (struct.is_audio_only != null) {
                protocol.mo14136v3("is_audio_only", 2, (byte) 2);
                protocol.mo14125P1(struct.is_audio_only.booleanValue());
            }
            protocol.mo14134i0();
        }
    }

    public AVCallMissed(@InterfaceC88465b Long l, @InterfaceC88465b Boolean bool) {
        this.sent_at_millis = l;
        this.is_audio_only = bool;
    }

    public static /* synthetic */ AVCallMissed copy$default(AVCallMissed aVCallMissed, Long l, Boolean bool, int i, Object obj) {
        if ((i & 1) != 0) {
            l = aVCallMissed.sent_at_millis;
        }
        if ((i & 2) != 0) {
            bool = aVCallMissed.is_audio_only;
        }
        return aVCallMissed.copy(l, bool);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final Long getSent_at_millis() {
        return this.sent_at_millis;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final Boolean getIs_audio_only() {
        return this.is_audio_only;
    }

    @InterfaceC88464a
    public final AVCallMissed copy(@InterfaceC88465b Long sent_at_millis, @InterfaceC88465b Boolean is_audio_only) {
        return new AVCallMissed(sent_at_millis, is_audio_only);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof AVCallMissed)) {
            return false;
        }
        AVCallMissed aVCallMissed = (AVCallMissed) other;
        return Intrinsics.m65267c(this.sent_at_millis, aVCallMissed.sent_at_millis) && Intrinsics.m65267c(this.is_audio_only, aVCallMissed.is_audio_only);
    }

    public int hashCode() {
        Long l = this.sent_at_millis;
        int iHashCode = (l == null ? 0 : l.hashCode()) * 31;
        Boolean bool = this.is_audio_only;
        return iHashCode + (bool != null ? bool.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "AVCallMissed(sent_at_millis=" + this.sent_at_millis + ", is_audio_only=" + this.is_audio_only + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
