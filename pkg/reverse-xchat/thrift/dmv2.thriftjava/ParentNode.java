package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import androidx.camera.camera2.internal.C0870z0;
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

@Metadata(m64929d1 = {"\u00004\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0002\b\u0004\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\b\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0006\b\u0086\b\u0018\u0000 \u001b2\u00020\u0001:\u0002\u001c\u001bB\u001b\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0004\u001a\u0004\u0018\u00010\u0002¢\u0006\u0004\b\u0005\u0010\u0006J\u0017\u0010\n\u001a\u00020\t2\u0006\u0010\b\u001a\u00020\u0007H\u0016¢\u0006\u0004\b\n\u0010\u000bJ\u0012\u0010\f\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\f\u0010\rJ\u0012\u0010\u000e\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000e\u0010\rJ(\u0010\u000f\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0004\u001a\u0004\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\u000f\u0010\u0010J\u0010\u0010\u0011\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0011\u0010\rJ\u0010\u0010\u0013\u001a\u00020\u0012HÖ\u0001¢\u0006\u0004\b\u0013\u0010\u0014J\u001a\u0010\u0018\u001a\u00020\u00172\b\u0010\u0016\u001a\u0004\u0018\u00010\u0015HÖ\u0003¢\u0006\u0004\b\u0018\u0010\u0019R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001aR\u0016\u0010\u0004\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001a¨\u0006\u001d"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ParentNode;", "Lcom/bendb/thrifty/a;", "", "subtree_encryption_public_key", "parent_hash", "<init>", "(Ljava/lang/String;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "copy", "(Ljava/lang/String;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/ParentNode;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Companion", "ParentNodeAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class ParentNode implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String parent_hash;

    @JvmField
    @InterfaceC88465b
    public final String subtree_encryption_public_key;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new ParentNodeAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ParentNode$ParentNodeAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/ParentNode;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/ParentNode;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/ParentNode;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class ParentNodeAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public ParentNode m83760read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            String string2 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new ParentNode(string, string2);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        C11272a.m14141a(protocol, b);
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

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a ParentNode struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("ParentNode");
            if (struct.subtree_encryption_public_key != null) {
                protocol.mo14136v3("subtree_encryption_public_key", 1, (byte) 11);
                protocol.mo14137w0(struct.subtree_encryption_public_key);
            }
            if (struct.parent_hash != null) {
                protocol.mo14136v3("parent_hash", 2, (byte) 11);
                protocol.mo14137w0(struct.parent_hash);
            }
            protocol.mo14134i0();
        }
    }

    public ParentNode(@InterfaceC88465b String str, @InterfaceC88465b String str2) {
        this.subtree_encryption_public_key = str;
        this.parent_hash = str2;
    }

    public static /* synthetic */ ParentNode copy$default(ParentNode parentNode, String str, String str2, int i, Object obj) {
        if ((i & 1) != 0) {
            str = parentNode.subtree_encryption_public_key;
        }
        if ((i & 2) != 0) {
            str2 = parentNode.parent_hash;
        }
        return parentNode.copy(str, str2);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getSubtree_encryption_public_key() {
        return this.subtree_encryption_public_key;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getParent_hash() {
        return this.parent_hash;
    }

    @InterfaceC88464a
    public final ParentNode copy(@InterfaceC88465b String subtree_encryption_public_key, @InterfaceC88465b String parent_hash) {
        return new ParentNode(subtree_encryption_public_key, parent_hash);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof ParentNode)) {
            return false;
        }
        ParentNode parentNode = (ParentNode) other;
        return Intrinsics.m65267c(this.subtree_encryption_public_key, parentNode.subtree_encryption_public_key) && Intrinsics.m65267c(this.parent_hash, parentNode.parent_hash);
    }

    public int hashCode() {
        String str = this.subtree_encryption_public_key;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        String str2 = this.parent_hash;
        return iHashCode + (str2 != null ? str2.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return C0870z0.m1255a("ParentNode(subtree_encryption_public_key=", this.subtree_encryption_public_key, ", parent_hash=", this.parent_hash, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
