package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import android.gov.nist.core.net.C0009a;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import com.google.ads.interactivemedia.p012v3.impl.data.C12781a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u00006\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0002\b\u0002\n\u0002\u0010\u000b\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\n\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0002\b\b\b\u0086\b\u0018\u0000 \u001f2\u00020\u0001:\u0002 \u001fB%\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0004\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0006\u001a\u0004\u0018\u00010\u0005¢\u0006\u0004\b\u0007\u0010\bJ\u0017\u0010\f\u001a\u00020\u000b2\u0006\u0010\n\u001a\u00020\tH\u0016¢\u0006\u0004\b\f\u0010\rJ\u0012\u0010\u000e\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000e\u0010\u000fJ\u0012\u0010\u0010\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u000fJ\u0012\u0010\u0011\u001a\u0004\u0018\u00010\u0005HÆ\u0003¢\u0006\u0004\b\u0011\u0010\u0012J4\u0010\u0013\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0004\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0006\u001a\u0004\u0018\u00010\u0005HÆ\u0001¢\u0006\u0004\b\u0013\u0010\u0014J\u0010\u0010\u0015\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u000fJ\u0010\u0010\u0017\u001a\u00020\u0016HÖ\u0001¢\u0006\u0004\b\u0017\u0010\u0018J\u001a\u0010\u001b\u001a\u00020\u00052\b\u0010\u001a\u001a\u0004\u0018\u00010\u0019HÖ\u0003¢\u0006\u0004\b\u001b\u0010\u001cR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001dR\u0016\u0010\u0004\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001dR\u0016\u0010\u0006\u001a\u0004\u0018\u00010\u00058\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010\u001e¨\u0006!"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/PullMessagePageDetails;", "Lcom/bendb/thrifty/a;", "", "min_sequence_id", "max_sequence_id", "", "is_batched_pull", "<init>", "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/Boolean;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "component3", "()Ljava/lang/Boolean;", "copy", "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/Boolean;)Lcom/x/dmv2/thriftjava/PullMessagePageDetails;", "toString", "", "hashCode", "()I", "", "other", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Ljava/lang/Boolean;", "Companion", "PullMessagePageDetailsAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class PullMessagePageDetails implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final Boolean is_batched_pull;

    @JvmField
    @InterfaceC88465b
    public final String max_sequence_id;

    @JvmField
    @InterfaceC88465b
    public final String min_sequence_id;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new PullMessagePageDetailsAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/PullMessagePageDetails$PullMessagePageDetailsAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/PullMessagePageDetails;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/PullMessagePageDetails;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/PullMessagePageDetails;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class PullMessagePageDetailsAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public PullMessagePageDetails m83765read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            String string2 = null;
            Boolean boolValueOf = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new PullMessagePageDetails(string, string2, boolValueOf);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 3) {
                    if (s != 4) {
                        if (s != 7) {
                            C11272a.m14141a(protocol, b);
                        } else if (b == 2) {
                            boolValueOf = Boolean.valueOf(protocol.readBool());
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

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a PullMessagePageDetails struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("PullMessagePageDetails");
            if (struct.min_sequence_id != null) {
                protocol.mo14136v3("min_sequence_id", 3, (byte) 11);
                protocol.mo14137w0(struct.min_sequence_id);
            }
            if (struct.max_sequence_id != null) {
                protocol.mo14136v3("max_sequence_id", 4, (byte) 11);
                protocol.mo14137w0(struct.max_sequence_id);
            }
            if (struct.is_batched_pull != null) {
                protocol.mo14136v3("is_batched_pull", 7, (byte) 2);
                protocol.mo14125P1(struct.is_batched_pull.booleanValue());
            }
            protocol.mo14134i0();
        }
    }

    public PullMessagePageDetails(@InterfaceC88465b String str, @InterfaceC88465b String str2, @InterfaceC88465b Boolean bool) {
        this.min_sequence_id = str;
        this.max_sequence_id = str2;
        this.is_batched_pull = bool;
    }

    public static /* synthetic */ PullMessagePageDetails copy$default(PullMessagePageDetails pullMessagePageDetails, String str, String str2, Boolean bool, int i, Object obj) {
        if ((i & 1) != 0) {
            str = pullMessagePageDetails.min_sequence_id;
        }
        if ((i & 2) != 0) {
            str2 = pullMessagePageDetails.max_sequence_id;
        }
        if ((i & 4) != 0) {
            bool = pullMessagePageDetails.is_batched_pull;
        }
        return pullMessagePageDetails.copy(str, str2, bool);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getMin_sequence_id() {
        return this.min_sequence_id;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getMax_sequence_id() {
        return this.max_sequence_id;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final Boolean getIs_batched_pull() {
        return this.is_batched_pull;
    }

    @InterfaceC88464a
    public final PullMessagePageDetails copy(@InterfaceC88465b String min_sequence_id, @InterfaceC88465b String max_sequence_id, @InterfaceC88465b Boolean is_batched_pull) {
        return new PullMessagePageDetails(min_sequence_id, max_sequence_id, is_batched_pull);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof PullMessagePageDetails)) {
            return false;
        }
        PullMessagePageDetails pullMessagePageDetails = (PullMessagePageDetails) other;
        return Intrinsics.m65267c(this.min_sequence_id, pullMessagePageDetails.min_sequence_id) && Intrinsics.m65267c(this.max_sequence_id, pullMessagePageDetails.max_sequence_id) && Intrinsics.m65267c(this.is_batched_pull, pullMessagePageDetails.is_batched_pull);
    }

    public int hashCode() {
        String str = this.min_sequence_id;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        String str2 = this.max_sequence_id;
        int iHashCode2 = (iHashCode + (str2 == null ? 0 : str2.hashCode())) * 31;
        Boolean bool = this.is_batched_pull;
        return iHashCode2 + (bool != null ? bool.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        String str = this.min_sequence_id;
        String str2 = this.max_sequence_id;
        return C12781a.m16257a(C0009a.m11b("PullMessagePageDetails(min_sequence_id=", str, ", max_sequence_id=", str2, ", is_batched_pull="), this.is_batched_pull, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}