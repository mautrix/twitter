package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
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

@Metadata(m64929d1 = {"\u00008\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010 \n\u0002\u0010\u000e\n\u0002\b\u0004\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\t\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000 \u001e2\u00020\u0001:\u0002\u001f\u001eB!\u0012\u000e\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0003¢\u0006\u0004\b\u0006\u0010\u0007J\u0017\u0010\u000b\u001a\u00020\n2\u0006\u0010\t\u001a\u00020\bH\u0016¢\u0006\u0004\b\u000b\u0010\fJ\u0018\u0010\r\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\r\u0010\u000eJ\u0012\u0010\u000f\u001a\u0004\u0018\u00010\u0003HÆ\u0003¢\u0006\u0004\b\u000f\u0010\u0010J.\u0010\u0011\u001a\u00020\u00002\u0010\b\u0002\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u0003HÆ\u0001¢\u0006\u0004\b\u0011\u0010\u0012J\u0010\u0010\u0013\u001a\u00020\u0003HÖ\u0001¢\u0006\u0004\b\u0013\u0010\u0010J\u0010\u0010\u0015\u001a\u00020\u0014HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u0016J\u001a\u0010\u001a\u001a\u00020\u00192\b\u0010\u0018\u001a\u0004\u0018\u00010\u0017HÖ\u0003¢\u0006\u0004\b\u001a\u0010\u001bR\u001c\u0010\u0004\u001a\n\u0012\u0004\u0012\u00020\u0003\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001cR\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00038\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010\u001d¨\u0006 "}, m64930d2 = {"Lcom/x/dmv2/thriftjava/UpdatePathNode;", "Lcom/bendb/thrifty/a;", "", "", "encrypted_secrets", "encrypted_private_key", "<init>", "(Ljava/util/List;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/util/List;", "component2", "()Ljava/lang/String;", "copy", "(Ljava/util/List;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/UpdatePathNode;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/util/List;", "Ljava/lang/String;", "Companion", "UpdatePathNodeAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class UpdatePathNode implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String encrypted_private_key;

    @JvmField
    @InterfaceC88465b
    public final List encrypted_secrets;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new UpdatePathNodeAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/UpdatePathNode$UpdatePathNodeAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/UpdatePathNode;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/UpdatePathNode;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/UpdatePathNode;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class UpdatePathNodeAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public UpdatePathNode m83789read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            ArrayList arrayList = null;
            String string = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new UpdatePathNode(arrayList, string);
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
                } else if (b == 15) {
                    int i = protocol.mo14130a2().f38395b;
                    ArrayList arrayList2 = new ArrayList(i);
                    for (int i2 = 0; i2 < i; i2++) {
                        arrayList2.add(protocol.readString());
                    }
                    arrayList = arrayList2;
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a UpdatePathNode struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("UpdatePathNode");
            if (struct.encrypted_secrets != null) {
                protocol.mo14136v3("encrypted_secrets", 1, (byte) 15);
                protocol.mo14128X0((byte) 11, struct.encrypted_secrets.size());
                Iterator it = struct.encrypted_secrets.iterator();
                while (it.hasNext()) {
                    protocol.mo14137w0((String) it.next());
                }
            }
            if (struct.encrypted_private_key != null) {
                protocol.mo14136v3("encrypted_private_key", 2, (byte) 11);
                protocol.mo14137w0(struct.encrypted_private_key);
            }
            protocol.mo14134i0();
        }
    }

    public UpdatePathNode(@InterfaceC88465b List list, @InterfaceC88465b String str) {
        this.encrypted_secrets = list;
        this.encrypted_private_key = str;
    }

    public static /* synthetic */ UpdatePathNode copy$default(UpdatePathNode updatePathNode, List list, String str, int i, Object obj) {
        if ((i & 1) != 0) {
            list = updatePathNode.encrypted_secrets;
        }
        if ((i & 2) != 0) {
            str = updatePathNode.encrypted_private_key;
        }
        return updatePathNode.copy(list, str);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final List getEncrypted_secrets() {
        return this.encrypted_secrets;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getEncrypted_private_key() {
        return this.encrypted_private_key;
    }

    @InterfaceC88464a
    public final UpdatePathNode copy(@InterfaceC88465b List encrypted_secrets, @InterfaceC88465b String encrypted_private_key) {
        return new UpdatePathNode(encrypted_secrets, encrypted_private_key);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof UpdatePathNode)) {
            return false;
        }
        UpdatePathNode updatePathNode = (UpdatePathNode) other;
        return Intrinsics.m65267c(this.encrypted_secrets, updatePathNode.encrypted_secrets) && Intrinsics.m65267c(this.encrypted_private_key, updatePathNode.encrypted_private_key);
    }

    public int hashCode() {
        List list = this.encrypted_secrets;
        int iHashCode = (list == null ? 0 : list.hashCode()) * 31;
        String str = this.encrypted_private_key;
        return iHashCode + (str != null ? str.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "UpdatePathNode(encrypted_secrets=" + this.encrypted_secrets + ", encrypted_private_key=" + this.encrypted_private_key + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
