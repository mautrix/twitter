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

@Metadata(m64929d1 = {"\u0000<\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\t\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\b\n\u0002\u0010\u000e\n\u0002\b\u0002\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0002\b\b\b\u0086\b\u0018\u0000 \u001f2\u00020\u0001:\u0002 \u001fB\u001b\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004¢\u0006\u0004\b\u0006\u0010\u0007J\u0017\u0010\u000b\u001a\u00020\n2\u0006\u0010\t\u001a\u00020\bH\u0016¢\u0006\u0004\b\u000b\u0010\fJ\u0012\u0010\r\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\r\u0010\u000eJ\u0012\u0010\u000f\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u000f\u0010\u0010J(\u0010\u0011\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u0004HÆ\u0001¢\u0006\u0004\b\u0011\u0010\u0012J\u0010\u0010\u0014\u001a\u00020\u0013HÖ\u0001¢\u0006\u0004\b\u0014\u0010\u0015J\u0010\u0010\u0017\u001a\u00020\u0016HÖ\u0001¢\u0006\u0004\b\u0017\u0010\u0018J\u001a\u0010\u001b\u001a\u00020\u00042\b\u0010\u001a\u001a\u0004\u0018\u00010\u0019HÖ\u0003¢\u0006\u0004\b\u001b\u0010\u001cR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001dR\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010\u001e¨\u0006!"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/SetVerifiedStatus;", "Lcom/bendb/thrifty/a;", "", "user_id", "", "verified_status", "<init>", "(Ljava/lang/Long;Ljava/lang/Boolean;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/Long;", "component2", "()Ljava/lang/Boolean;", "copy", "(Ljava/lang/Long;Ljava/lang/Boolean;)Lcom/x/dmv2/thriftjava/SetVerifiedStatus;", "", "toString", "()Ljava/lang/String;", "", "hashCode", "()I", "", "other", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/Long;", "Ljava/lang/Boolean;", "Companion", "SetVerifiedStatusAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class SetVerifiedStatus implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final Long user_id;

    @JvmField
    @InterfaceC88465b
    public final Boolean verified_status;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new SetVerifiedStatusAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/SetVerifiedStatus$SetVerifiedStatusAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/SetVerifiedStatus;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/SetVerifiedStatus;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/SetVerifiedStatus;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class SetVerifiedStatusAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public SetVerifiedStatus m85671read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Long lValueOf = null;
            Boolean boolValueOf = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new SetVerifiedStatus(lValueOf, boolValueOf);
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

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a SetVerifiedStatus struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("SetVerifiedStatus");
            if (struct.user_id != null) {
                protocol.mo14136v3("user_id", 1, (byte) 10);
                protocol.mo14121B3(struct.user_id.longValue());
            }
            if (struct.verified_status != null) {
                protocol.mo14136v3("verified_status", 2, (byte) 2);
                protocol.mo14125P1(struct.verified_status.booleanValue());
            }
            protocol.mo14134i0();
        }
    }

    public SetVerifiedStatus(@InterfaceC88465b Long l, @InterfaceC88465b Boolean bool) {
        this.user_id = l;
        this.verified_status = bool;
    }

    public static /* synthetic */ SetVerifiedStatus copy$default(SetVerifiedStatus setVerifiedStatus, Long l, Boolean bool, int i, Object obj) {
        if ((i & 1) != 0) {
            l = setVerifiedStatus.user_id;
        }
        if ((i & 2) != 0) {
            bool = setVerifiedStatus.verified_status;
        }
        return setVerifiedStatus.copy(l, bool);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final Long getUser_id() {
        return this.user_id;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final Boolean getVerified_status() {
        return this.verified_status;
    }

    @InterfaceC88464a
    public final SetVerifiedStatus copy(@InterfaceC88465b Long user_id, @InterfaceC88465b Boolean verified_status) {
        return new SetVerifiedStatus(user_id, verified_status);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof SetVerifiedStatus)) {
            return false;
        }
        SetVerifiedStatus setVerifiedStatus = (SetVerifiedStatus) other;
        return Intrinsics.m65267c(this.user_id, setVerifiedStatus.user_id) && Intrinsics.m65267c(this.verified_status, setVerifiedStatus.verified_status);
    }

    public int hashCode() {
        Long l = this.user_id;
        int iHashCode = (l == null ? 0 : l.hashCode()) * 31;
        Boolean bool = this.verified_status;
        return iHashCode + (bool != null ? bool.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "SetVerifiedStatus(user_id=" + this.user_id + ", verified_status=" + this.verified_status + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}