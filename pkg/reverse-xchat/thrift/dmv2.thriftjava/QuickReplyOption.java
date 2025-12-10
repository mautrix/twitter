package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import android.gov.nist.core.net.C0009a;
import android.gov.nist.javax.sdp.fields.C0015d;
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

@Metadata(m64929d1 = {"\u00004\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0002\b\u0006\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\n\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0006\b\u0086\b\u0018\u0000 \u001f2\u00020\u0001:\u0002 \u001fB/\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0004\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0006\u001a\u0004\u0018\u00010\u0002¢\u0006\u0004\b\u0007\u0010\bJ\u0017\u0010\f\u001a\u00020\u000b2\u0006\u0010\n\u001a\u00020\tH\u0016¢\u0006\u0004\b\f\u0010\rJ\u0012\u0010\u000e\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000e\u0010\u000fJ\u0012\u0010\u0010\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u000fJ\u0012\u0010\u0011\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0011\u0010\u000fJ\u0012\u0010\u0012\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0012\u0010\u000fJ@\u0010\u0013\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0004\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0006\u001a\u0004\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\u0013\u0010\u0014J\u0010\u0010\u0015\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u000fJ\u0010\u0010\u0017\u001a\u00020\u0016HÖ\u0001¢\u0006\u0004\b\u0017\u0010\u0018J\u001a\u0010\u001c\u001a\u00020\u001b2\b\u0010\u001a\u001a\u0004\u0018\u00010\u0019HÖ\u0003¢\u0006\u0004\b\u001c\u0010\u001dR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001eR\u0016\u0010\u0004\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001eR\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010\u001eR\u0016\u0010\u0006\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010\u001e¨\u0006!"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReplyOption;", "Lcom/bendb/thrifty/a;", "", "id", "label", "metadata", "description", "<init>", "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "component3", "component4", "copy", "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/QuickReplyOption;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Companion", "QuickReplyOptionAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class QuickReplyOption implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String description;

    @JvmField
    @InterfaceC88465b
    public final String id;

    @JvmField
    @InterfaceC88465b
    public final String label;

    @JvmField
    @InterfaceC88465b
    public final String metadata;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new QuickReplyOptionAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReplyOption$QuickReplyOptionAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/QuickReplyOption;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/QuickReplyOption;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/QuickReplyOption;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class QuickReplyOptionAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public QuickReplyOption m83769read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            String string2 = null;
            String string3 = null;
            String string4 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new QuickReplyOption(string, string2, string3, string4);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            if (s != 4) {
                                C11272a.m14141a(protocol, b);
                            } else if (b == 11) {
                                string4 = protocol.readString();
                            } else {
                                C11272a.m14141a(protocol, b);
                            }
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

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a QuickReplyOption struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("QuickReplyOption");
            if (struct.id != null) {
                protocol.mo14136v3("id", 1, (byte) 11);
                protocol.mo14137w0(struct.id);
            }
            if (struct.label != null) {
                protocol.mo14136v3("label", 2, (byte) 11);
                protocol.mo14137w0(struct.label);
            }
            if (struct.metadata != null) {
                protocol.mo14136v3("metadata", 3, (byte) 11);
                protocol.mo14137w0(struct.metadata);
            }
            if (struct.description != null) {
                protocol.mo14136v3("description", 4, (byte) 11);
                protocol.mo14137w0(struct.description);
            }
            protocol.mo14134i0();
        }
    }

    public QuickReplyOption(@InterfaceC88465b String str, @InterfaceC88465b String str2, @InterfaceC88465b String str3, @InterfaceC88465b String str4) {
        this.id = str;
        this.label = str2;
        this.metadata = str3;
        this.description = str4;
    }

    public static /* synthetic */ QuickReplyOption copy$default(QuickReplyOption quickReplyOption, String str, String str2, String str3, String str4, int i, Object obj) {
        if ((i & 1) != 0) {
            str = quickReplyOption.id;
        }
        if ((i & 2) != 0) {
            str2 = quickReplyOption.label;
        }
        if ((i & 4) != 0) {
            str3 = quickReplyOption.metadata;
        }
        if ((i & 8) != 0) {
            str4 = quickReplyOption.description;
        }
        return quickReplyOption.copy(str, str2, str3, str4);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getId() {
        return this.id;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getLabel() {
        return this.label;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final String getMetadata() {
        return this.metadata;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final String getDescription() {
        return this.description;
    }

    @InterfaceC88464a
    public final QuickReplyOption copy(@InterfaceC88465b String id, @InterfaceC88465b String label, @InterfaceC88465b String metadata, @InterfaceC88465b String description) {
        return new QuickReplyOption(id, label, metadata, description);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof QuickReplyOption)) {
            return false;
        }
        QuickReplyOption quickReplyOption = (QuickReplyOption) other;
        return Intrinsics.m65267c(this.id, quickReplyOption.id) && Intrinsics.m65267c(this.label, quickReplyOption.label) && Intrinsics.m65267c(this.metadata, quickReplyOption.metadata) && Intrinsics.m65267c(this.description, quickReplyOption.description);
    }

    public int hashCode() {
        String str = this.id;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        String str2 = this.label;
        int iHashCode2 = (iHashCode + (str2 == null ? 0 : str2.hashCode())) * 31;
        String str3 = this.metadata;
        int iHashCode3 = (iHashCode2 + (str3 == null ? 0 : str3.hashCode())) * 31;
        String str4 = this.description;
        return iHashCode3 + (str4 != null ? str4.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        String str = this.id;
        String str2 = this.label;
        return C0015d.m22a(C0009a.m11b("QuickReplyOption(id=", str, ", label=", str2, ", metadata="), this.metadata, ", description=", this.description, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}