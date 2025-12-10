package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import androidx.media3.exoplayer.mediacodec.C8338z;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000>\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0000\n\u0002\u0010 \n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\t\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000 \u001f2\u00020\u0001:\u0002 \u001fB!\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\u000e\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u0004¢\u0006\u0004\b\u0007\u0010\bJ\u0017\u0010\f\u001a\u00020\u000b2\u0006\u0010\n\u001a\u00020\tH\u0016¢\u0006\u0004\b\f\u0010\rJ\u0012\u0010\u000e\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000e\u0010\u000fJ\u0018\u0010\u0010\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u0011J.\u0010\u0012\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\u0010\b\u0002\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u0004HÆ\u0001¢\u0006\u0004\b\u0012\u0010\u0013J\u0010\u0010\u0014\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0014\u0010\u000fJ\u0010\u0010\u0016\u001a\u00020\u0015HÖ\u0001¢\u0006\u0004\b\u0016\u0010\u0017J\u001a\u0010\u001b\u001a\u00020\u001a2\b\u0010\u0019\u001a\u0004\u0018\u00010\u0018HÖ\u0003¢\u0006\u0004\b\u001b\u0010\u001cR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001dR\u001c\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010\u001e¨\u0006!"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReplyOptionsRequest;", "Lcom/bendb/thrifty/a;", "", "id", "", "Lcom/x/dmv2/thriftjava/QuickReplyOption;", "options", "<init>", "(Ljava/lang/String;Ljava/util/List;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "()Ljava/util/List;", "copy", "(Ljava/lang/String;Ljava/util/List;)Lcom/x/dmv2/thriftjava/QuickReplyOptionsRequest;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Ljava/util/List;", "Companion", "QuickReplyOptionsRequestAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class QuickReplyOptionsRequest implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String id;

    @JvmField
    @InterfaceC88465b
    public final List options;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new QuickReplyOptionsRequestAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReplyOptionsRequest$QuickReplyOptionsRequestAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/QuickReplyOptionsRequest;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/QuickReplyOptionsRequest;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/QuickReplyOptionsRequest;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class QuickReplyOptionsRequestAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public QuickReplyOptionsRequest m83770read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            ArrayList arrayList = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new QuickReplyOptionsRequest(string, arrayList);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        C11272a.m14141a(protocol, b);
                    } else if (b == 15) {
                        int i = protocol.mo14130a2().f38395b;
                        ArrayList arrayList2 = new ArrayList(i);
                        for (int i2 = 0; i2 < i; i2++) {
                            arrayList2.add((QuickReplyOption) QuickReplyOption.ADAPTER.read(protocol));
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

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a QuickReplyOptionsRequest struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("QuickReplyOptionsRequest");
            if (struct.id != null) {
                protocol.mo14136v3("id", 1, (byte) 11);
                protocol.mo14137w0(struct.id);
            }
            if (struct.options != null) {
                protocol.mo14136v3("options", 2, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.options.size());
                Iterator it = struct.options.iterator();
                while (it.hasNext()) {
                    QuickReplyOption.ADAPTER.write(protocol, (QuickReplyOption) it.next());
                }
            }
            protocol.mo14134i0();
        }
    }

    public QuickReplyOptionsRequest(@InterfaceC88465b String str, @InterfaceC88465b List list) {
        this.id = str;
        this.options = list;
    }

    public static /* synthetic */ QuickReplyOptionsRequest copy$default(QuickReplyOptionsRequest quickReplyOptionsRequest, String str, List list, int i, Object obj) {
        if ((i & 1) != 0) {
            str = quickReplyOptionsRequest.id;
        }
        if ((i & 2) != 0) {
            list = quickReplyOptionsRequest.options;
        }
        return quickReplyOptionsRequest.copy(str, list);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getId() {
        return this.id;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final List getOptions() {
        return this.options;
    }

    @InterfaceC88464a
    public final QuickReplyOptionsRequest copy(@InterfaceC88465b String id, @InterfaceC88465b List options) {
        return new QuickReplyOptionsRequest(id, options);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof QuickReplyOptionsRequest)) {
            return false;
        }
        QuickReplyOptionsRequest quickReplyOptionsRequest = (QuickReplyOptionsRequest) other;
        return Intrinsics.m65267c(this.id, quickReplyOptionsRequest.id) && Intrinsics.m65267c(this.options, quickReplyOptionsRequest.options);
    }

    public int hashCode() {
        String str = this.id;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        List list = this.options;
        return iHashCode + (list != null ? list.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return C8338z.m10466b("QuickReplyOptionsRequest(id=", this.id, ", options=", this.options, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}