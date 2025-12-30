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

@Metadata(m64929d1 = {"\u00004\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\t\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0002\b\b\b\u0086\b\u0018\u0000 \u001d2\u00020\u0001:\u0002\u001e\u001dB\u001b\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004¢\u0006\u0004\b\u0006\u0010\u0007J\u0017\u0010\u000b\u001a\u00020\n2\u0006\u0010\t\u001a\u00020\bH\u0016¢\u0006\u0004\b\u000b\u0010\fJ\u0012\u0010\r\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\r\u0010\u000eJ\u0012\u0010\u000f\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u000f\u0010\u0010J(\u0010\u0011\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u0004HÆ\u0001¢\u0006\u0004\b\u0011\u0010\u0012J\u0010\u0010\u0013\u001a\u00020\u0004HÖ\u0001¢\u0006\u0004\b\u0013\u0010\u0010J\u0010\u0010\u0015\u001a\u00020\u0014HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u0016J\u001a\u0010\u0019\u001a\u00020\u00022\b\u0010\u0018\u001a\u0004\u0018\u00010\u0017HÖ\u0003¢\u0006\u0004\b\u0019\u0010\u001aR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001bR\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010\u001c¨\u0006\u001f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/AVCallStarted;", "Lcom/bendb/thrifty/a;", "", "is_audio_only", "", "broadcast_id", "<init>", "(Ljava/lang/Boolean;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/Boolean;", "component2", "()Ljava/lang/String;", "copy", "(Ljava/lang/Boolean;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/AVCallStarted;", "toString", "", "hashCode", "()I", "", "other", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/Boolean;", "Ljava/lang/String;", "Companion", "AVCallStartedAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class AVCallStarted implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String broadcast_id;

    @JvmField
    @InterfaceC88465b
    public final Boolean is_audio_only;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new AVCallStartedAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/AVCallStarted$AVCallStartedAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/AVCallStarted;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/AVCallStarted;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/AVCallStarted;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class AVCallStartedAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public AVCallStarted m85638read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Boolean boolValueOf = null;
            String string = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new AVCallStarted(boolValueOf, string);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 3) {
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
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a AVCallStarted struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("AVCallStarted");
            if (struct.is_audio_only != null) {
                protocol.mo14136v3("is_audio_only", 1, (byte) 2);
                protocol.mo14125P1(struct.is_audio_only.booleanValue());
            }
            if (struct.broadcast_id != null) {
                protocol.mo14136v3("broadcast_id", 3, (byte) 11);
                protocol.mo14137w0(struct.broadcast_id);
            }
            protocol.mo14134i0();
        }
    }

    public AVCallStarted(@InterfaceC88465b Boolean bool, @InterfaceC88465b String str) {
        this.is_audio_only = bool;
        this.broadcast_id = str;
    }

    public static /* synthetic */ AVCallStarted copy$default(AVCallStarted aVCallStarted, Boolean bool, String str, int i, Object obj) {
        if ((i & 1) != 0) {
            bool = aVCallStarted.is_audio_only;
        }
        if ((i & 2) != 0) {
            str = aVCallStarted.broadcast_id;
        }
        return aVCallStarted.copy(bool, str);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final Boolean getIs_audio_only() {
        return this.is_audio_only;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getBroadcast_id() {
        return this.broadcast_id;
    }

    @InterfaceC88464a
    public final AVCallStarted copy(@InterfaceC88465b Boolean is_audio_only, @InterfaceC88465b String broadcast_id) {
        return new AVCallStarted(is_audio_only, broadcast_id);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof AVCallStarted)) {
            return false;
        }
        AVCallStarted aVCallStarted = (AVCallStarted) other;
        return Intrinsics.m65267c(this.is_audio_only, aVCallStarted.is_audio_only) && Intrinsics.m65267c(this.broadcast_id, aVCallStarted.broadcast_id);
    }

    public int hashCode() {
        Boolean bool = this.is_audio_only;
        int iHashCode = (bool == null ? 0 : bool.hashCode()) * 31;
        String str = this.broadcast_id;
        return iHashCode + (str != null ? str.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "AVCallStarted(is_audio_only=" + this.is_audio_only + ", broadcast_id=" + this.broadcast_id + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
