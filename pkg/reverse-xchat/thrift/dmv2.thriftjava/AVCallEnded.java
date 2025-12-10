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

@Metadata(m64929d1 = {"\u0000<\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\t\n\u0002\b\u0002\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\f\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0002\b\t\b\u0086\b\u0018\u0000 $2\u00020\u0001:\u0002%$B/\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0004\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0006\u001a\u0004\u0018\u00010\u0005\u0012\b\u0010\b\u001a\u0004\u0018\u00010\u0007¢\u0006\u0004\b\t\u0010\nJ\u0017\u0010\u000e\u001a\u00020\r2\u0006\u0010\f\u001a\u00020\u000bH\u0016¢\u0006\u0004\b\u000e\u0010\u000fJ\u0012\u0010\u0010\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u0011J\u0012\u0010\u0012\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0012\u0010\u0011J\u0012\u0010\u0013\u001a\u0004\u0018\u00010\u0005HÆ\u0003¢\u0006\u0004\b\u0013\u0010\u0014J\u0012\u0010\u0015\u001a\u0004\u0018\u00010\u0007HÆ\u0003¢\u0006\u0004\b\u0015\u0010\u0016J@\u0010\u0017\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0004\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0006\u001a\u0004\u0018\u00010\u00052\n\b\u0002\u0010\b\u001a\u0004\u0018\u00010\u0007HÆ\u0001¢\u0006\u0004\b\u0017\u0010\u0018J\u0010\u0010\u0019\u001a\u00020\u0007HÖ\u0001¢\u0006\u0004\b\u0019\u0010\u0016J\u0010\u0010\u001b\u001a\u00020\u001aHÖ\u0001¢\u0006\u0004\b\u001b\u0010\u001cJ\u001a\u0010\u001f\u001a\u00020\u00052\b\u0010\u001e\u001a\u0004\u0018\u00010\u001dHÖ\u0003¢\u0006\u0004\b\u001f\u0010 R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010!R\u0016\u0010\u0004\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010!R\u0016\u0010\u0006\u001a\u0004\u0018\u00010\u00058\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010\"R\u0016\u0010\b\u001a\u0004\u0018\u00010\u00078\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\b\u0010#¨\u0006&"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/AVCallEnded;", "Lcom/bendb/thrifty/a;", "", "sent_at_millis", "duration_seconds", "", "is_audio_only", "", "broadcast_id", "<init>", "(Ljava/lang/Long;Ljava/lang/Long;Ljava/lang/Boolean;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/Long;", "component2", "component3", "()Ljava/lang/Boolean;", "component4", "()Ljava/lang/String;", "copy", "(Ljava/lang/Long;Ljava/lang/Long;Ljava/lang/Boolean;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/AVCallEnded;", "toString", "", "hashCode", "()I", "", "other", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/Long;", "Ljava/lang/Boolean;", "Ljava/lang/String;", "Companion", "AVCallEndedAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class AVCallEnded implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String broadcast_id;

    @JvmField
    @InterfaceC88465b
    public final Long duration_seconds;

    @JvmField
    @InterfaceC88465b
    public final Boolean is_audio_only;

    @JvmField
    @InterfaceC88465b
    public final Long sent_at_millis;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new AVCallEndedAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/AVCallEnded$AVCallEndedAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/AVCallEnded;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/AVCallEnded;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/AVCallEnded;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class AVCallEndedAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public AVCallEnded m85636read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Long lValueOf = null;
            Long lValueOf2 = null;
            Boolean boolValueOf = null;
            String string = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new AVCallEnded(lValueOf, lValueOf2, boolValueOf, string);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            if (s != 5) {
                                C11272a.m14141a(protocol, b);
                            } else if (b == 11) {
                                string = protocol.readString();
                            } else {
                                C11272a.m14141a(protocol, b);
                            }
                        } else if (b == 2) {
                            boolValueOf = Boolean.valueOf(protocol.readBool());
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 10) {
                        lValueOf2 = Long.valueOf(protocol.mo14124H0());
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

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a AVCallEnded struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("AVCallEnded");
            if (struct.sent_at_millis != null) {
                protocol.mo14136v3("sent_at_millis", 1, (byte) 10);
                protocol.mo14121B3(struct.sent_at_millis.longValue());
            }
            if (struct.duration_seconds != null) {
                protocol.mo14136v3("duration_seconds", 2, (byte) 10);
                protocol.mo14121B3(struct.duration_seconds.longValue());
            }
            if (struct.is_audio_only != null) {
                protocol.mo14136v3("is_audio_only", 3, (byte) 2);
                protocol.mo14125P1(struct.is_audio_only.booleanValue());
            }
            if (struct.broadcast_id != null) {
                protocol.mo14136v3("broadcast_id", 5, (byte) 11);
                protocol.mo14137w0(struct.broadcast_id);
            }
            protocol.mo14134i0();
        }
    }

    public AVCallEnded(@InterfaceC88465b Long l, @InterfaceC88465b Long l2, @InterfaceC88465b Boolean bool, @InterfaceC88465b String str) {
        this.sent_at_millis = l;
        this.duration_seconds = l2;
        this.is_audio_only = bool;
        this.broadcast_id = str;
    }

    public static /* synthetic */ AVCallEnded copy$default(AVCallEnded aVCallEnded, Long l, Long l2, Boolean bool, String str, int i, Object obj) {
        if ((i & 1) != 0) {
            l = aVCallEnded.sent_at_millis;
        }
        if ((i & 2) != 0) {
            l2 = aVCallEnded.duration_seconds;
        }
        if ((i & 4) != 0) {
            bool = aVCallEnded.is_audio_only;
        }
        if ((i & 8) != 0) {
            str = aVCallEnded.broadcast_id;
        }
        return aVCallEnded.copy(l, l2, bool, str);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final Long getSent_at_millis() {
        return this.sent_at_millis;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final Long getDuration_seconds() {
        return this.duration_seconds;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final Boolean getIs_audio_only() {
        return this.is_audio_only;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final String getBroadcast_id() {
        return this.broadcast_id;
    }

    @InterfaceC88464a
    public final AVCallEnded copy(@InterfaceC88465b Long sent_at_millis, @InterfaceC88465b Long duration_seconds, @InterfaceC88465b Boolean is_audio_only, @InterfaceC88465b String broadcast_id) {
        return new AVCallEnded(sent_at_millis, duration_seconds, is_audio_only, broadcast_id);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof AVCallEnded)) {
            return false;
        }
        AVCallEnded aVCallEnded = (AVCallEnded) other;
        return Intrinsics.m65267c(this.sent_at_millis, aVCallEnded.sent_at_millis) && Intrinsics.m65267c(this.duration_seconds, aVCallEnded.duration_seconds) && Intrinsics.m65267c(this.is_audio_only, aVCallEnded.is_audio_only) && Intrinsics.m65267c(this.broadcast_id, aVCallEnded.broadcast_id);
    }

    public int hashCode() {
        Long l = this.sent_at_millis;
        int iHashCode = (l == null ? 0 : l.hashCode()) * 31;
        Long l2 = this.duration_seconds;
        int iHashCode2 = (iHashCode + (l2 == null ? 0 : l2.hashCode())) * 31;
        Boolean bool = this.is_audio_only;
        int iHashCode3 = (iHashCode2 + (bool == null ? 0 : bool.hashCode())) * 31;
        String str = this.broadcast_id;
        return iHashCode3 + (str != null ? str.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "AVCallEnded(sent_at_millis=" + this.sent_at_millis + ", duration_seconds=" + this.duration_seconds + ", is_audio_only=" + this.is_audio_only + ", broadcast_id=" + this.broadcast_id + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}